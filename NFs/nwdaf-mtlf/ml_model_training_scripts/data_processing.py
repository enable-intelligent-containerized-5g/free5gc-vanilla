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

def combine_csv_files(base_file, additional_file, output_file, cpu_colum, mem_column):
    fieldnames = ['namespace', 'pod', 'container', 'timestamp', cpu_colum, mem_column]
    
    # Open the base file
    if os.path.isfile(base_file):
        with open(base_file, mode='r') as basefile:
            reader_base = csv.DictReader(basefile)
            base_data = list(reader_base)
    else:
        with open(base_file, mode='w') as basefile:
            writer = csv.DictWriter(basefile, fieldnames=fieldnames)
            # Write the header with the new column
            writer.writeheader()
            base_data = []  # Empty array
    
    # Open the aditional file
    with open(additional_file, mode='r') as additionalfile:
        reader_additional = csv.DictReader(additionalfile)
        additional_data_dict = {row['timestamp']: row[mem_column] for row in reader_additional}

    # Create or open a CSV file to save the combined results
    with open(output_file, mode='w', newline='') as combinedfile:
        writer = csv.DictWriter(combinedfile, fieldnames=fieldnames)
        # Write the header with the new column
        writer.writeheader()
        # Iterate over the base file and add the additional column
        for row in base_data:
            timestamp = row['timestamp']
            # Add value
            row[mem_column] = additional_data_dict.get(timestamp, 'N/A')  # If there is no value, assign 'N/A'

            # Write the row
            writer.writerow(row)
            
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
    if len(sys.argv) < 14:
        sys.exit("Missing params")

    # Getting params
    loc_script = sys.argv[0]
    data_path = sys.argv[1]
    data_raw_path = sys.argv[2]
    data_preprocessed_path = sys.argv[3]
    data_processed_path = sys.argv[4]
    data_labeled_path = sys.argv[5]
    cpu_file = sys.argv[6]
    mem_file = sys.argv[7]
    thrpt_file = sys.argv[8]
    dataset_file = sys.argv[9]
    selected_dataset_file = sys.argv[10]
    cpu_column = sys.argv[11]
    mem_column = sys.argv[12]
    thrpt_column = sys.argv[13]
    
    # Validate folders
    isFolder([data_path, data_raw_path, data_preprocessed_path, data_processed_path, data_labeled_path])

    # Define the input and output data
    cpu_file_path = data_raw_path + cpu_file
    mem_file_path = data_raw_path + mem_file
    thrpt_file_path = data_raw_path + thrpt_file
    intermediate_cpu_csv = data_preprocessed_path + "cpuUsage.csv"
    intermediate_mem_csv = data_preprocessed_path + "memUsage.csv"
    intermediate_thrpt_csv = data_preprocessed_path + "thrptTotal.csv"
    intermediate_data_csv = data_preprocessed_path + "combinedData.csv"
    output_cpu_csv = data_processed_path + "processedCpuUsage.csv"
    output_mem_csv = data_processed_path + "processedMemUsage.csv"
    output_thrpt_csv = data_processed_path + "processedThrptTotal.csv"
    output_data_csv = data_labeled_path + dataset_file
    selected_dataset_path = data_labeled_path + selected_dataset_file
    
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
    
    process_json_to_csv(thrpt_file_path, intermediate_thrpt_csv, thrpt_column)
    clean_csv(intermediate_thrpt_csv, output_thrpt_csv)
    
    # Fieldnames
    fieldnames_base = ['namespace', 'pod', 'container', 'timestamp', cpu_column, mem_column]
    fieldnames_cpu = ['namespace', 'pod', 'container', 'timestamp', cpu_column]
    fieldnames_mem = ['namespace', 'pod', 'container', 'timestamp', mem_column]
    fieldnames_thrpt = ['namespace', 'pod', 'container', 'timestamp', thrpt_column]
    fieldnames_common = ['namespace', 'pod', 'container', 'timestamp']

    # Validate csv
    df_output = validateCsv(output_data_csv, fieldnames_base)
    df_base = validateCsv(selected_dataset_path, fieldnames_base)
    df_cpu = validateCsv(output_cpu_csv, fieldnames_cpu)
    df_mem = validateCsv(output_mem_csv, fieldnames_mem)
    df_thrpt = validateCsv(output_thrpt_csv, fieldnames_thrpt)
        
    # Perform merge based on key columns
    merged_cpu_mem = pd.merge(df_cpu, df_mem[fieldnames_mem],
                        on=fieldnames_common,
                        how='left')
    # Perform merge based on key columns
    merged_resources_thrpt = pd.merge(merged_cpu_mem, df_thrpt[fieldnames_thrpt],
                        on=fieldnames_common,
                        how='left')
    
    
    # Perform a merge to find rows in merged_resources_thrpt that are not in df_base
    merged = pd.merge(df_base[fieldnames_common], merged_resources_thrpt, 
                    on=fieldnames_common, 
                    how='right', indicator=True)
    # Filter out rows in merged_resources_thrpt that are not in df_base
    cpu_mem_missing = merged[merged['_merge'] == 'right_only'].drop(columns=['_merge'])
    # Add missing rows from merged_resources_thrpt to df_base
    df_final = pd.concat([df_base, cpu_mem_missing], ignore_index=True)
    
    
    # Sort the DataFrame by 'timestamp' in ascending order
    df_sorted = df_final.sort_values(by='timestamp', ascending=True)
    # Delete rows with NaN columns
    df_clean = df_sorted.dropna()
    # Save dataset
    df_output = df_clean
    # exit(f"{output_data_csv}")
    df_output.to_csv(output_data_csv, index=False)
        
    # combine_csv_files(output_cpu_csv, output_mem_csv, intermediate_data_csv, cpu_column, mem_column)
    # clean_csv(intermediate_data_csv, output_data_csv)

if __name__ == "__main__":
    main()

