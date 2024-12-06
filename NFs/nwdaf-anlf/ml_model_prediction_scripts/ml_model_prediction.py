import sys, csv, joblib, os, json
import numpy as np
import pandas as pd
from sklearn.preprocessing import MinMaxScaler

def result_validate(data, cpu_pred, mem_pred):
    # Must be a numpy array
    if not isinstance(data, np.ndarray):
        return False, f"The prediction has an incorrect format"

    # Must be two-dimensional
    if data.ndim != 2:
        return False, f"The result does not have the predictions of {cpu_pred} and {mem_pred}"

    # Check that it has exactly 1 row and 2 columns
    if data.shape[1] != 2:
        return False, f"The result does not have 2 predictions, has {data.shape[1]}"

    return True, "The prediction has a correct format"


def nf_load_prediction(data_path, dataset_path, prediction_file , cpu_column, mem_column, selected_model_uri, time_steps):
    
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
    
    if len(data_values) != time_steps :
        sys.exit(f"The dataset does not have the required number of rows, provides a larger dataset")
    
    # Scale the data between 0 and 1
    scaler = MinMaxScaler(feature_range=(0, 1))
    data_scaled = scaler.fit_transform(data_values) # Comun Dataset

    
    
    ##################################################################
    ###           Load Ramdom Forest Regresos Ml Model             ###
    ##################################################################
    
    try:
        ml_model = joblib.load(selected_model_uri)
    except FileNotFoundError as e:
        sys.exit(f"No Found the Ml Model: {e}")
    except Exception as e:
        sys.exit(f"Error loading the Ml Model: {e}")
    
    ##################################################################
    ###                          Prediction                        ###
    ##################################################################

    # Info
    new_input = data_scaled.reshape(1, -1)
       
    # Make the predictions
    y_pred = ml_model.predict(new_input)
    
    # Invert the normalization to obtain the original values
    y_pred_invertido = scaler.inverse_transform(y_pred)
    
    result, err_pred = result_validate(y_pred_invertido, cpu_column, mem_column)
    if result == False :
        sys.exit(err_pred)
        
    cpu_prediction = y_pred_invertido[0][0]
    mem_prediction = y_pred_invertido[0][1]
    
    
    ##################################################################
    ###                    Save prediction info                    ###
    ##################################################################
        
    # Save model
    prediction_info_path = f"{data_path}{prediction_file}"
    
    prediction_info = {
        'cpuAverage': cpu_prediction,
        'memAverage': mem_prediction
    }
    
    with open(prediction_info_path, 'w') as json_file:
        json.dump(prediction_info, json_file, indent=4)                           
        
        
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
    if len(sys.argv) < 10:
        sys.exit(1)

    # Get the params
    loc_script = sys.argv[0]
    models_path = sys.argv[1]
    data_path = sys.argv[2]
    data_labeled_path = sys.argv[3]
    dataset_file = sys.argv[4]
    prediction_file = sys.argv[5]
    cpu_column = sys.argv[6]
    mem_column = sys.argv[7]
    selected_model_uri = sys.argv[8]
    time_steps = sys.argv[9]
    
    # Validate time_steps
    try:
        int_time_steps = int(time_steps.strip())
    except ValueError:
        sys.exit(f"Invalid input: not a valid timeSteps value")
    
    # Validate folders
    isFolder([models_path, data_path, data_labeled_path])
    
    dataset_path = data_labeled_path + dataset_file
    
    # Validate dataset_path
    isFile([dataset_path, selected_model_uri])

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
        
    nf_load_prediction(data_path, dataset_path, prediction_file, cpu_column, mem_column, selected_model_uri, int_time_steps)

if __name__ == "__main__":
    main()
