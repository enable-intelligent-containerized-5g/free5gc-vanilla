import sys
import json
import csv
import pandas as pd

def process_json_to_csv(input_file, output_file, colum_value):
    # Leer el archivo JSON
    with open(input_file, 'r') as file:
        data = json.load(file)

    # Verificar que el archivo no esté vacío
    if not data:
        print("El archivo JSON está vacío.")
        sys.exit(1)

    # Abrir el archivo CSV para escritura
    with open(output_file, 'w', newline='') as csvfile:
        fieldnames = ['namespace', 'pod', 'container', 
                      'timestamp', colum_value]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        # Escribir encabezados
        writer.writeheader()

        # Procesar los datos
        for i in range(len(data)):
            # Tomar los últimos 4 valores de 'value'
            row = {
                "namespace": data[i]['namespace'],
                "pod": data[i]['pod'],
                "container": data[i]['container'],
                "timestamp": data[i]['timestamp'],
                colum_value: data[i]['value']
            }
            writer.writerow(row)

    # print(f"Archivo CSV generado en: {output_file}")
    
def clean_csv(input_csv, output_csv):
    # Leer el archivo CSV en un DataFrame de pandas
    df = pd.read_csv(input_csv)

    # Eliminar filas con valores nulos
    df.dropna(inplace=True)

    # Guardar el DataFrame limpio en un nuevo archivo CSV
    df.to_csv(output_csv, index=False)
    # print(f"CSV file created in: {output_csv}")

def combine_csv_files(base_file, additional_file, output_file, cpu_colum, mem_column):
    # Abrir el archivo base
    with open(base_file, mode='r') as basefile:
        reader_base = csv.DictReader(basefile)
        base_data = list(reader_base)  # Leer todo el archivo base en memoria

    # Abrir el archivo adicional
    with open(additional_file, mode='r') as additionalfile:
        reader_additional = csv.DictReader(additionalfile)
        additional_data = {row['timestamp']: row[mem_column] for row in reader_additional}  # Crear un diccionario con timestamp como clave

    # Crear o abrir un archivo CSV para guardar los resultados combinados
    with open(output_file, mode='w', newline='') as combinedfile:
        fieldnames = ['namespace', 'pod', 'container', 'timestamp', cpu_colum, mem_column]
        writer = csv.DictWriter(combinedfile, fieldnames=fieldnames)

        # Escribir la cabecera con la nueva columna 'additional_value'
        writer.writeheader()

        # Iterar sobre el archivo base y agregar la columna adicional 'additional_value'
        for row in base_data:
            timestamp = row['timestamp']
            # Si el timestamp coincide en el archivo adicional, agregar el valor
            row[mem_column] = additional_data.get(timestamp, 'N/A')  # Si no hay valor, asignar 'N/A'

            # Escribir la fila combinada en el archivo de salida
            writer.writerow(row)

def main():
    # Verificar si se pasaron los parámetros correctos
    if len(sys.argv) < 4:
        print("Error: Not enough parameters received.", file=sys.stderr)
        sys.exit(1)

    # Obtener los parámetros desde la línea de comandos
    loc_script = sys.argv[0]
    base_path = sys.argv[1]
    cpu_file = sys.argv[2]
    mem_file = sys.argv[3]
    print(f"\nReceived parameters:\n- Base Path: {base_path}\n- CPU File: {cpu_file}\n- Mem File: {mem_file}")
    
    cpu_path = base_path + cpu_file
    mem_path = base_path + mem_file

    # Intentar cargar los archivos JSON
    try:
        with open(cpu_path, 'r') as f_cpu:
            cpu_json = json.load(f_cpu)
        with open(mem_path, 'r') as f_mem:
            mem_json = json.load(f_mem)

        print("Files uploaded successfully.")
    except FileNotFoundError as e:
        print(f"Error opening files: {e}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}", file=sys.stderr)
        sys.exit(1)

    # Definir los archivos de entrada y salida
    intermediate_cpu_csv = base_path + "cpuUsage.csv"
    intermediate_mem_csv = base_path + "memUage.csv"
    intermediate_data_csv = base_path + "data.csv"
    output_cpu_csv = base_path + "processedCpuUsage.csv"
    output_mem_csv = base_path + "processedMemUsage.csv"
    output_data_csv = base_path + "processedData.csv"
    cpu_column = "cpu_usage"
    mem_column = "mem_usage"

    # Ejecutar la función
    process_json_to_csv(cpu_path, intermediate_cpu_csv, cpu_column)
    clean_csv(intermediate_cpu_csv, output_cpu_csv)
    
    process_json_to_csv(mem_path, intermediate_mem_csv, mem_column)
    clean_csv(intermediate_mem_csv, output_mem_csv)
    
    combine_csv_files(output_cpu_csv, output_mem_csv, intermediate_data_csv, cpu_column, mem_column)
    clean_csv(intermediate_data_csv, output_data_csv)
    
    print(f"Data processing completed and saved in: {output_data_csv}")

if __name__ == "__main__":
    main()

