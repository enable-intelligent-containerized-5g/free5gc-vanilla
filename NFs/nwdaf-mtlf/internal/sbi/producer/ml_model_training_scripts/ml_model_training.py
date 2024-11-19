import sys, csv, joblib, os, json
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler
from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import mean_squared_error, r2_score


def ml_model_training(models_path, data_path, figures_path, dataset_path, model_info_file, model_info_list, cpu_column, mem_column, base_name, time_steps):
    
    ##################################################################
    ###                 Get and Process the data                   ###
    ##################################################################
    
    # Load data from a CSV file
    def load_data_from_csv(csv_file):
        data = pd.read_csv(csv_file)
        return data

    # Load dataset from a CSV file
    df = load_data_from_csv(dataset_path)
    
    # We select the columns that we are going to use for the prediction
    data_values = df[[cpu_column, mem_column]].values
    
    # Validate the dataset size
    # time_steps = 4 # Steps
    test_size = 0.3 # testing part
    if (len(data_values)-time_steps)*test_size < 1:
        sys.exit(f"The dataset does not have the required number of rows, provides a larger dataset")
        
    # Scale the data between 0 and 1
    scaler = MinMaxScaler(feature_range=(0, 1))
    data_scaled = scaler.fit_transform(data_values) # Comun Dataset
    
    # Función para crear las secuencias
    def create_sequences_multivariate(data, time_steps):
        X, y = [], []
        for i in range(len(data) - time_steps):
            X.append(data[i:i + time_steps])  # Seleccionamos las últimas 'time_steps' filas (como secuencia)
            y.append(data[i + time_steps])  # Valores para predecir
        return np.array(X), np.array(y)
    
    X, y = create_sequences_multivariate(data_scaled, time_steps)
    
    
    ##################################################################
    ###        Random Forest Regressor model configuration         ###
    ##################################################################

    X_train, X_test, y_train, y_test = train_test_split(X.reshape(X.shape[0], -1), y, test_size=test_size, random_state=42)

    # Create the model
    rf_model = RandomForestRegressor(n_estimators=100, random_state=42)

    # Train the model
    rf_model.fit(X_train, y_train)
    
    ##################################################################
    ###                          Evaluation                        ###
    ##################################################################

    # Info
    name = "RFR"
    large_name = 'Random Forest Regressor'
       
    # Make the predictions
    y_pred = rf_model.predict(X_test)
    # Invert the normalization to obtain the original values
    y_pred_invertido = scaler.inverse_transform(y_pred)
    y_test_invertido = scaler.inverse_transform(y_test)
    
    # Evaluate the model
    mse = mean_squared_error(y_test_invertido, y_pred_invertido)
    r2 = r2_score(y_test_invertido, y_pred_invertido)
    metrics = f'MSE: {mse:.4f}, R²: {r2:.4f}'
    
    # Evaluate the model: MSE and R² for each output (CPU and Memory)
    mse_cpu = mean_squared_error(y_test_invertido[:, 0], y_pred_invertido[:, 0])  # For CPU column
    mse_mem = mean_squared_error(y_test_invertido[:, 1], y_pred_invertido[:, 1])  # for Memory column
    r2_cpu = r2_score(y_test_invertido[:, 0], y_pred_invertido[:, 0])  # R² for CPU
    r2_mem = r2_score(y_test_invertido[:, 1], y_pred_invertido[:, 1])  # R² fot Memory
    cpu_metrics = f'CPU - R²: {r2_cpu:.4f}, MSE: {mse_cpu:.4f}'
    mem_metrics = f'Memory - R²: {r2_mem:.4f}, MSE: {mse_mem:.4f}'
    
    
    ##################################################################
    ###                      Plot the results                      ###
    ##################################################################
    
    # Create the figure
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 6))

    # CPU Graph
    ax1.scatter(y_test_invertido[:, 0], y_pred_invertido[:, 0], color='blue', label='Prediction vs Real CPU')
    ax1.plot([min(y_test_invertido[:, 0]), max(y_test_invertido[:, 0])], 
            [min(y_test_invertido[:, 0]), max(y_test_invertido[:, 0])], color='red', linestyle='--', label='CPU reference line')
    ax1.set_xlabel('Real CPU Usage')
    ax1.set_ylabel('Predicted CPU Usage')
    ax1.set_title(f'CPU Predictions (MSE: {mse_cpu:.4f}, R²: {r2_cpu:.4f})')
    ax1.legend()
    ax1.grid(True)
    # Memory graph
    ax2.scatter(y_test_invertido[:, 1], y_pred_invertido[:, 1], color='green', label='Prediction vs Real Memory')
    ax2.plot([min(y_test_invertido[:, 1]), max(y_test_invertido[:, 1])], 
            [min(y_test_invertido[:, 1]), max(y_test_invertido[:, 1])], color='orange', linestyle='--', label='Memory reference line')
    ax2.set_xlabel('Real Memory Usage')
    ax2.set_ylabel('Predicted Memory Usage')
    ax2.set_title(f'Memory Predictions (MSE: {mse_mem:.4f}, R²: {r2_mem:.4f})')
    ax2.legend()
    ax2.grid(True)

    # Title
    fig.suptitle(f'{large_name} ({name}) model\nMSE: {mse:.4f}, R²: {r2:.4f}', fontsize=14)
    # Adjust the graphs
    plt.tight_layout(pad=0.8) 
    # Show plot
    # plt.show()
    
    # Save the plot
    fig_format = "png"
    
    fig_name = f"model_{name}_{base_name}"
    fig_uri = f"{figures_path}{fig_name}.{fig_format}"
    plt.savefig(fig_uri, format=fig_format, bbox_inches='tight')
    
    
    ##################################################################
    ###                       Save model info                      ###
    ##################################################################
        
    # Save model
    model_format = 'pkl'
    model_info_path = f"{data_path}{model_info_file}"
    
    model_name = fig_name
    model_uri = f"{models_path}{model_name}.{model_format}"
    joblib.dump(rf_model, model_uri)
    model_size = os.path.getsize(model_uri)
    
    model_info = {
        'name': model_name,
        'uri': model_uri,
        'size': model_size,
        'figureUri': fig_uri,
        'mse': mse,
        'r2':r2,
        'mseCpu': mse_cpu,
        'r2Cpu':r2_cpu,
        'mseMem': mse_mem,
        'r2Mem':r2_mem,
    }
    
    with open(model_info_path, 'w') as json_file:
        json.dump(model_info, json_file, indent=4)
        
    # Add model_info to list of model_info_list
    model_info_list_path = data_path + model_info_list
    models_info = []
    try:
        with open(model_info_list_path, 'r') as json_file:
            content = json_file.read()
            models_info = json.loads(content)
            if content.strip(): # No empty
                if isinstance(models_info, list):
                    models_info = json.loads(content)
                       
               
    except FileNotFoundError:
        models_info = []
        print("No found the models info list. Creating a new list")
    except json.JSONDecodeError as e:
        models_info = []
        print("Error decoding the models info file. Creating a new list")

    models_info.append(model_info)

    with open(model_info_list_path, 'w') as json_file:
        json.dump(models_info, json_file, indent=4)
        
        
def isFile(file_paths):
    for file_path in file_paths:
        if not os.path.isfile(file_path):
            sys.exit(f"The specified path '{file_path}' is not a valid file.")
        
def isFolder(folder_paths):
    for folder_path in folder_paths:
        if not os.path.isdir(folder_path):
            sys.exit(f"The specified path '{folder_path}' is not a valid directory.")
       

def main():
    # Verify the params
    if len(sys.argv) < 12:
        sys.exit(1)

    # Get the params
    loc_script = sys.argv[0]
    models_path = sys.argv[1]
    data_path = sys.argv[2]
    data_labeled_path = sys.argv[3]
    figures_path = sys.argv[4]
    dataset_file = sys.argv[5]
    model_info = sys.argv[6]
    model_info_list = sys.argv[7]
    cpu_column = sys.argv[8]
    mem_column = sys.argv[9]
    base_name = sys.argv[10]
    time_steps = sys.argv[11]
    
    # Validate folders
    isFolder([models_path, data_path, data_labeled_path, figures_path])
    
    dataset_path = data_labeled_path + dataset_file
    
    # Validate dataset_path
    isFile([dataset_path])
    
    # Validate time_steps
    try:
        int_time_steps = int(time_steps.strip())
    except ValueError:
        sys.exit(f"Invalid input: not a valid timeSteps value")

    # Try load the data 
    try:
        # Try to open the dataset file
        with open(dataset_path, mode='r') as file:
            reader = csv.DictReader(file)
            
    except FileNotFoundError:
        sys.exit(f"No found the dataset {dataset_file}")
    except csv.Error as e:
        sys.exit(f"Error opening the dataset {dataset_file}")
    except Exception as e:
        sys.exit(f"Error opening the dataset {dataset_file}")
        
    ml_model_training(models_path, data_path, figures_path, dataset_path, model_info, model_info_list, cpu_column, mem_column, base_name, int_time_steps)

if __name__ == "__main__":
    main()
