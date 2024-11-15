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
        return

    # Abrir el archivo CSV para escritura
    with open(output_file, 'w', newline='') as csvfile:
        fieldnames = ['namespace', 'pod', 'container', 
                      'timestamp', 'cpu_usage']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        # Escribir encabezados
        writer.writeheader()

        cpu_limit = 0.25
        # Procesar los datos
        for i in range(3, len(data)):
            # Tomar los últimos 4 valores de 'value'
            row = {
                "namespace": data[i]['namespace'],
                "pod": data[i]['pod'],
                "container": data[i]['container'],
                "timestamp": data[i]['timestamp'],
                "cpu_usage": data[i]['value']/cpu_limit
            }
            writer.writerow(row)

    print(f"Archivo CSV generado en: {output_file}")
    
def clean_csv(input_csv, output_csv):
    # Leer el archivo CSV en un DataFrame de pandas
    df = pd.read_csv(input_csv)

    # Eliminar filas con valores nulos
    df.dropna(inplace=True)

    # Guardar el DataFrame limpio en un nuevo archivo CSV
    df.to_csv(output_csv, index=False)
    print(f"Archivo CSV limpio generado en: {output_csv}")


# Definir los archivos de entrada y salida
input_file = 'output5.json'
intermediate_csv = 'lsmt_data3.csv'
output_file = 'processed_lsmt_data3.csv'

# Ejecutar la función
process_json_to_csv(input_file, intermediate_csv)
clean_csv(intermediate_csv, output_file)
