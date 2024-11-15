import sys
import json
import csv
import pandas as pd

def process_json_to_csv(input_file, output_file):
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
                      'timestamp', 'cpu_usage']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        # Escribir encabezados
        writer.writeheader()

        # Procesar los datos
        for i in range(3, len(data)):
            # Tomar los últimos 4 valores de 'value'
            row = {
                "namespace": data[i]['namespace'],
                "pod": data[i]['pod'],
                "container": data[i]['container'],
                "timestamp": data[i]['timestamp'],
                "cpu_usage": data[i]['value']
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
    print(f"CSV file created in: {output_csv}")

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
    print(f"Received parameters:\n- Base Path: {base_path}\n- CPU File: {cpu_file}\n- Mem File: {mem_file}")
    
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
    output_cpu_csv = base_path + "cpuUsageProcessed.csv"
    output_mem_csv = base_path + "memUsageProcessed.csv"

    # Ejecutar la función
    process_json_to_csv(cpu_path, intermediate_cpu_csv)
    clean_csv(intermediate_cpu_csv, output_cpu_csv)
    
    process_json_to_csv(mem_path, intermediate_mem_csv)
    clean_csv(intermediate_mem_csv, output_mem_csv)

if __name__ == "__main__":
    main()

