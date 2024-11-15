# train_model.py

import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sklearn.metrics import mean_squared_error, r2_score
from xgboost import XGBRegressor
from sklearn.tree import DecisionTreeRegressor
from sklearn.ensemble import RandomForestRegressor


# Cargar datos desde un archivo CSV
def load_data_from_csv(csv_file):
    data = pd.read_csv(csv_file)
    return data

# Cargar el conjunto de datos desde un archivo CSV
path_data = 'processed_lsmt_data3.csv'
df = load_data_from_csv(path_data)

# Definir el tamaño de la secuencia
time_steps = 3

# Función para crear las secuencias
def create_sequences_multivariate(data, time_steps):
    X, y = [], []
    for i in range(len(data) - time_steps):
        X.append(data[i:i + time_steps])  # Seleccionamos las últimas 'time_steps' filas (como secuencia)
        y.append(data[i + time_steps, 0])  # El valor de la columna CPU_Usage para predecir
    return np.array(X), np.array(y)

# Seleccionamos las columnas que vamos a usar para la predicción
data_values = df[['cpu_usage']].values

# Crear las secuencias
X, y = create_sequences_multivariate(data_values, time_steps)

X_train, X_test, y_train, y_test = train_test_split(X.reshape(X.shape[0], -1), y, test_size=0.3, random_state=42)

# Crear el modelo de XGBoost
xgb_model = XGBRegressor(n_estimators=100, random_state=42)
rf_model = RandomForestRegressor(n_estimators=100, random_state=42)

# Entrenar el modelo
xgb_model.fit(X_train, y_train)
rf_model.fit(X_train, y_train)

#Evaluate the models
for model, name in zip([xgb_model, rf_model], ['XGBoost', 'rf_model']):
    print(f"MODEL: {name}")
    # Predecir en el conjunto de prueba

    # Realizar predicciones
    y_pred = model.predict(X_test)

    # Evaluar el modelo
    mse = mean_squared_error(y_test, y_pred)
    r2 = r2_score(y_test, y_pred)
    print(f'MSE: {mse:.4f}')
    print(f'R²: {r2:.4f}')
    print()
    
    
