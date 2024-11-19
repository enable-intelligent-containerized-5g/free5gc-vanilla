import sys, json, csv, os
import pandas as pd
import numpy as np

def process_json_to_csv(input_file, output_file, colum_value):
    # Read JSON file
    with open(input_file, 'r') as file:
        data = json.load(file)

    # Check that the file is not empty
    if not data:
        sys.exit(1)

    # Open CSV file for writing
    with open(output_file, 'w', newline='') as csvfile:
        fieldnames = ['namespace', 'pod', 'container', 
                      'timestamp', colum_value]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        # Writing Headlines
        writer.writeheader()

        # Process the data
        for i in range(len(data)):
            # Take the last 4 values ​​of 'value'
            row = {
                "namespace": data[i]['namespace'],
                "pod": data[i]['pod'],
                "container": data[i]['container'],
                "timestamp": data[i]['timestamp'],
                colum_value: data[i]['value']
            }
            writer.writerow(row)
    
def clean_csv(input_csv, output_csv):
    # Read CSV file into a pandas DataFrame
    df = pd.read_csv(input_csv)

    # Delete rows with null values
    df.dropna(inplace=True)

    # Save the cleaned DataFrame to a new CSV file
    df.to_csv(output_csv, index=False)
    
def isFile(file_paths):
    for file_path in file_paths:
        if not os.path.isfile(file_path):
            sys.exit(f"The specified path '{file_path}' is not a valid file.")
        
def isFolder(folder_paths):
    for folder_path in folder_paths:
        if not os.path.isdir(folder_path):
            sys.exit(f"The specified path '{folder_path}' is not a valid directory.")

            
def validateCsv(file_csv, fieldnames):
    if os.path.exists(file_csv):
        if os.path.getsize(file_csv) > 0:
            try:
                df = pd.read_csv(file_csv)
                missing_columns = set(fieldnames) - set(df.columns)
                if missing_columns:
                    for col in missing_columns:
                        df[col] = pd.NA
            except pd.errors.EmptyDataError:
                df = pd.DataFrame(columns=fieldnames)
        else:
            df = pd.DataFrame(columns=fieldnames)
    else:
        df = pd.DataFrame(columns=fieldnames)
    
    return df

def main():
    # Verify the params
    if len(sys.argv) < 11:
        sys.exit(1)

    # Getting params
    loc_script = sys.argv[0]
    data_path = sys.argv[1]
    data_raw_path = sys.argv[2]
    data_preprocessed_path = sys.argv[3]
    data_processed_path = sys.argv[4]
    data_labeled_path = sys.argv[5]
    cpu_file = sys.argv[6]
    mem_file = sys.argv[7]
    dataset_file = sys.argv[8]
    cpu_column = sys.argv[9]
    mem_column = sys.argv[10]
    
    # Validate folders
    isFolder([data_path, data_raw_path, data_preprocessed_path, data_processed_path, data_labeled_path])

    # Define the input and output data
    cpu_file_path = data_raw_path + cpu_file
    mem_file_path = data_raw_path + mem_file
    intermediate_cpu_csv = data_preprocessed_path + "cpuUsage.csv"
    intermediate_mem_csv = data_preprocessed_path + "memUsage.csv"
    output_cpu_csv = data_processed_path + "processedCpuUsage.csv"
    output_mem_csv = data_processed_path + "processedMemUsage.csv"
    output_data_csv = data_labeled_path + dataset_file
    
    # Validate Files
    isFile([cpu_file_path, mem_file_path])

    # Trying to load the JSON files
    try:
        with open(cpu_file_path, 'r') as f_cpu:
            cpu_json = json.load(f_cpu)
        with open(mem_file_path, 'r') as f_mem:
            mem_json = json.load(f_mem)

    except FileNotFoundError as e:
        sys.exit(f"{e}")
    except json.JSONDecodeError as e:
        sys.exit(f"{e}")

    # Process the data
    process_json_to_csv(cpu_file_path, intermediate_cpu_csv, cpu_column)
    clean_csv(intermediate_cpu_csv, output_cpu_csv)
    
    process_json_to_csv(mem_file_path, intermediate_mem_csv, mem_column)
    clean_csv(intermediate_mem_csv, output_mem_csv)
    
    # Fieldnames
    fieldnames_cpu = ['namespace', 'pod', 'container', 'timestamp', cpu_column]
    fieldnames_mem = ['namespace', 'pod', 'container', 'timestamp', mem_column]
    fieldnames_common = ['namespace', 'pod', 'container', 'timestamp']

    # Validate csv
    df_cpu = validateCsv(output_cpu_csv, fieldnames_cpu)
    df_mem = validateCsv(output_mem_csv, fieldnames_mem)
        
    # Perform merge based on key columns
    merged_cpu_mem = pd.merge(df_cpu, df_mem[fieldnames_mem],
                        on=fieldnames_common,
                        how='left')
    df_sorted = merged_cpu_mem.sort_values(by='timestamp', ascending=True)
    # Delete rows with NaN columns
    df_clean = df_sorted.dropna()
    # Save dataset
    df_clean.to_csv(output_data_csv, index=False)

if __name__ == "__main__":
    main()

