import time, datetime, joblib, os, json
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler
from xgboost import XGBRegressor
from sklearn.tree import DecisionTreeRegressor
from sklearn.ensemble import RandomForestRegressor
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, GRU, Dense
from collections import namedtuple


def switch_case(fig_name, model, model_type, models_path):
    if model_type == 'xgboost':
        format = '.json'
        name = f"{fig_name}{format}"
        uri = f"{models_path}{name}"
        model.save_model(uri)
        size = os.path.getsize(uri)
        return uri, size
    
    elif model_type == 'sklearn':
        format = '.pkl'
        name = f"{fig_name}{format}"
        uri = f"{models_path}{name}"
        joblib.dump(model, uri)
        size = os.path.getsize(uri)
        return uri, size
    
    elif model_type == 'keras':
        format = '.h5'
        name = f"{fig_name}{format}"
        uri = f"{models_path}{name}"
        model.save(uri)
        size = os.path.getsize(uri)
        return uri, size
    
    else:
        return "none", 0

def plot_results(y_test_invertido, y_pred_invertido, name, large_name, model, model_type):
    # Evaluate the model
    mse = mean_squared_error(y_test_invertido, y_pred_invertido)
    r2 = r2_score(y_test_invertido, y_pred_invertido)
    print(f'MSE: {mse:.4f}, R²: {r2:.4f}')
    
    # Evaluate the model: MSE and R² for each output (CPU and Memory)
    mse_cpu = mean_squared_error(y_test_invertido[:, 0], y_pred_invertido[:, 0])  # Para la columna de CPU
    mse_mem = mean_squared_error(y_test_invertido[:, 1], y_pred_invertido[:, 1])  # Para la columna de Memoria
    r2_cpu = r2_score(y_test_invertido[:, 0], y_pred_invertido[:, 0])  # R² para CPU
    r2_mem = r2_score(y_test_invertido[:, 1], y_pred_invertido[:, 1])  # R² para Memoria
    
    print(f'CPU - R²: {r2_cpu:.4f}, MSE: {mse_cpu:.4f}')
    print(f'Memory - R²: {r2_mem:.4f}, MSE: {mse_mem:.4f}')
    
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
    
    # Save plot
    fig_path = "figures/"
    fig_format = "svg"
    fig_name, fig_uri = save_figure(plt, fig_path, name, fig_format)
    
    # Save model
    models_path = "saved_models/"
    info_models_path = "models_info.json"
    model_uri, size = switch_case(fig_name, model, model_type, models_path)
    
    if model_uri != "none":
        # Save info
        new_model = {
            'name': name,
            'uri': model_uri,
            'size': size,
            'figure': fig_uri,
            'mse': mse,
            'r2':r2,
            'mse_cpu': mse_cpu,
            'r2_cpu':r2_cpu,
            'mse_mem': mse_mem,
            'r2_mem':r2_mem,
        }
        
        try:
            with open(info_models_path, 'r') as json_file:
                models_info = json.load(json_file)
        except FileNotFoundError:
            models_info = [] 

        models_info.append(new_model)

        with open(info_models_path, 'w') as json_file:
            json.dump(models_info, json_file, indent=4)
    
    # Plot the results for CPU and Memory
    # plt.figure(figsize=(12, 6))
    # plt.scatter(y_test_invertido[:, 0], y_pred_invertido[:, 0], color='blue', label='Prediction vs Real CPU')  # Para CPU
    # plt.scatter(y_test_invertido[:, 1], y_pred_invertido[:, 1], color='green', label='Prediction vs Real Memory')  # Para Memoria
    # plt.plot([min(y_test_invertido[:, 0]), max(y_test_invertido[:, 0])], [min(y_test_invertido[:, 0]), max(y_test_invertido[:, 0])], color='red', linestyle='--', label='CPU reference line')
    # plt.plot([min(y_test_invertido[:, 1]), max(y_test_invertido[:, 1])], [min(y_test_invertido[:, 1]), max(y_test_invertido[:, 1])], color='orange', linestyle='--', label='Memory reference line')
    # # Add tags and title
    # plt.xlabel('Real Values')
    # plt.ylabel('Predictions')
    # plt.title(f'{large_name} ({name}) (MSE: {mse:.4f}, R²: {r2:.4f})\nMSE(CPU) = {mse_cpu:.3f}, MSE(Mem) = {mse_mem:.3f}\nR²(CPU) = {r2_cpu:.3f}, R²(Mem) = {r2_mem:.3f}')
    # plt.legend()
    # Show the graph
    # plt.show()

    
def save_figure(plot, fig_path, model_name, format):
    time.sleep(1)
    current_date = datetime.datetime.now()
    formated_current_date = current_date.strftime("%Y-%m-%d") + "_" + current_date.strftime("%H-%M-%S") + "_" + str(current_date.microsecond)
    
    fig_name = f"{model_name}_{formated_current_date}"
    fig_uri = f"{fig_path}{fig_name}.{format}"
    plot.savefig(fig_uri, format=format, bbox_inches='tight')
    
    return fig_name, fig_uri


def ml_model_training(data_path, cpu_column, mem_column):
    ##################################################################
    ###                   Common configuration                     ###
    ##################################################################
    
    # Load data from a CSV file
    def load_data_from_csv(csv_file):
        data = pd.read_csv(csv_file)
        return data

    # Load dataset from a CSV file
    df = load_data_from_csv(data_path)
    
    # We select the columns that we are going to use for the prediction
    data_values = df[[cpu_column, mem_column]].values
    # Scale the data between 0 and 1
    scaler = MinMaxScaler(feature_range=(0, 1))
    data_scaled = scaler.fit_transform(data_values) # Comun Dataset
    time_steps = 4 # Steps
    
    # Función para crear las secuencias
    def create_sequences_multivariate(data, time_steps):
        X, y = [], []
        for i in range(len(data) - time_steps):
            X.append(data[i:i + time_steps])  # Seleccionamos las últimas 'time_steps' filas (como secuencia)
            y.append(data[i + time_steps])  # Valores para predecir
        return np.array(X), np.array(y)
    
    X, y = create_sequences_multivariate(data_scaled, time_steps)
    
    
    
    ##################################################################
    ###                         LSTM, GRU                          ###
    ##################################################################

    if False :
        # Dividir los datos en entrenamiento (70%) y prueba (30%) de forma secuencial
        train_size = int(len(X) * 0.7)
        X_train, X_test = X[:train_size], X[train_size:]
        y_train, y_test = y[:train_size], y[train_size:]

        # Define the LSTM model
        lstm_model = Sequential()
        lstm_model.add(LSTM(100, return_sequences=True, input_shape=(time_steps, X.shape[2])))
        lstm_model.add(LSTM(50))
        lstm_model.add(Dense(2))
        lstm_model.compile(optimizer='adam', loss='mse')
        #  Train the model
        history = lstm_model.fit(X_train, y_train, epochs=30, batch_size=32, validation_data=(X_test, y_test))
        
        # Defining the GRU model
        gru_model = Sequential()
        gru_model.add(GRU(100, return_sequences=True, input_shape=(time_steps, X.shape[2])))
        gru_model.add(GRU(50))
        gru_model.add(Dense(2)) 
        gru_model.compile(optimizer='adam', loss='mse')
        # Train the model
        history = gru_model.fit(X_train, y_train, epochs=30, batch_size=32, validation_data=(X_test, y_test))
        
        #Evaluate the models
        for model, name, large_name, model_type in zip([lstm_model, gru_model], ['LSTM', 'GRU'], ['Long Short-Term Memory', 'Gated Recurrent Unit'], ['keras', 'keras']):
            print()
            print(f"MODEL: {large_name}")

            # Make predictions
            y_pred = model.predict(X_test)
            # Invert the normalization to obtain the original values
            y_pred_invertido = scaler.inverse_transform(y_pred)
            y_test_invertido = scaler.inverse_transform(y_test)
            
            # Plot
            plot_results(y_test_invertido, y_pred_invertido, name, large_name, model, model_type)

            # # Graficar la pérdida durante el entrenamiento
            # plt.figure(figsize=(10, 6))
            # plt.plot(history.history['loss'], label='Pérdida de Entrenamiento')
            # plt.plot(history.history['val_loss'], label='Pérdida de Validación')
            # plt.title('Pérdida durante el Entrenamiento')
            # plt.xlabel('Épocas')
            # plt.ylabel('MSE')
            # plt.legend()
            # plt.show()

    
    
    ##################################################################
    ### XGBRegressor, RandomForestRegressor, DecisionTreeRegressor ###
    ###                    and LinearRegression                    ###
    ##################################################################

    if True :
        X_train, X_test, y_train, y_test = train_test_split(X.reshape(X.shape[0], -1), y, test_size=0.3, random_state=42)

        # Create the models
        xgb_model = XGBRegressor(n_estimators=100, random_state=42)
        rf_model = RandomForestRegressor(n_estimators=100, random_state=42)
        dt_model = DecisionTreeRegressor(random_state=42)
        lr_model = LinearRegression()

        # Train the model
        xgb_model.fit(X_train, y_train)
        rf_model.fit(X_train, y_train)
        dt_model.fit(X_train, y_train)
        lr_model.fit(X_train, y_train)

        # Evaluate the models
        for model, name, large_name, model_type in zip([xgb_model, rf_model, dt_model, lr_model], ['XGBoost', 'RF', 'DT', 'LR'], ['Extreme Gradient Boosting', 'Random Forest', 'Decision Tree', 'Linear Regression'], ['xgboost', 'sklearn', 'sklearn', 'sklearn']):
            print()
            print(f"MODEL: {large_name}")

            # Make the predictions
            y_pred = model.predict(X_test)
            # Invert the normalization to obtain the original values
            y_pred_invertido = scaler.inverse_transform(y_pred)
            y_test_invertido = scaler.inverse_transform(y_test)
            
            # Plot
            plot_results(y_test_invertido, y_pred_invertido, name, large_name, model, model_type)
            
        
        
    ##################################################################
    ###                    Multilayer Perceptron                   ###
    ##################################################################
            
    # Create the lags features
    if False :
        def create_lagged_features(data, lag):
            X, y = [], []
            for i in range(lag, len(data)):
                X.append(data[i-lag:i].flatten())
                y.append(data[i])
            return np.array(X), np.array(y)
        
        lag = time_steps
        X, y = create_lagged_features(data_scaled, lag)
        
        # Divide the data
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.3, random_state=42)
        
        # Define the MLP model
        mlp_model = Sequential()
        mlp_model.add(Dense(64, activation='relu', input_shape=(X_train.shape[1],)))
        mlp_model.add(Dense(32, activation='relu'))
        mlp_model.add(Dense(2))  # Salida para predecir tanto CPU como Memoria
        mlp_model.compile(optimizer='adam', loss='mse')

        # Train the model
        history = mlp_model.fit(X_train, y_train, epochs=30, batch_size=32, validation_split=0.2)

        print()
        name = "MLP"
        model_type = 'keras'
        large_name = "Multilayer Perceptron"
        print(f"MODEL: {name}")
        # Make the predictions
        y_pred = mlp_model.predict(X_test)
        
        # Invert teh data to get the real values
        y_pred_invertido = scaler.inverse_transform(y_pred)
        y_test_invertido = scaler.inverse_transform(y_test)
        
        # Plot
        plot_results(y_test_invertido, y_pred_invertido, name, large_name, mlp_model, model_type)

def main():
    print("Ml Model Training")
    
    # Params
    data_path = "dataset.csv"
    cpu_column = "cpu_usage"
    mem_column = "mem_usage"
        
    ml_model_training(data_path, cpu_column, mem_column)

if __name__ == "__main__":
    main()
    

    
    
