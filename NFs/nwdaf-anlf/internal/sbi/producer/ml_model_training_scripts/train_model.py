# train_model.py

import numpy as np
import pandas as pd
from xgboost import XGBRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.tree import DecisionTreeRegressor
from sklearn.ensemble import RandomForestRegressor

# Cargar datos desde un archivo CSV
def load_data_from_csv(csv_file):
    data = pd.read_csv(csv_file)
    return data

# Cargar el conjunto de datos desde un archivo CSV
path_data = 'processed_xgboost4.csv'
data = load_data_from_csv(path_data)

features = data[['cpu_usage_1', 'cpu_usage_2', 'cpu_usage_3']]  # Aquí seleccionas las características que deseas usar
target = data['cpu_usage_4']  # La columna objetivo
# print(features.var())
# print(features.corr())

# Normalizar las características
scalerMinMax = MinMaxScaler(feature_range=(0, 1))  # Escalar entre 0 y 1
scalerStandard = StandardScaler()
features_scaled = scalerStandard.fit_transform(features)
features_scaled = scalerMinMax.fit_transform(features)

# print(features)

# Dividir el conjunto de datos en entrenamiento y prueba
X_train, X_test, y_train, y_test = train_test_split(features, target, test_size=0.3, random_state=42)

# Inicializar el modelo XGBRegressor
#Initialize models
dt_model = DecisionTreeRegressor(random_state=42)
rf_model = RandomForestRegressor(n_estimators=100, random_state=42)
xgb_model = XGBRegressor(objective='reg:squarederror', eval_metric='rmse')


# Entrenar el modelo
dt_model.fit(X_train, y_train)
rf_model.fit(X_train, y_train)
xgb_model.fit(X_train, y_train)

#Evaluate the models
for model, name in zip([dt_model, rf_model, xgb_model], ['Decision Tree', 'Random Forest', 'XGBoost']):
    print()
    print(f"{name} Model:")
    # Predecir en el conjunto de prueba
    y_pred = model.predict(X_test)

    # Evaluar el modelo
    mse = mean_squared_error(y_test, y_pred)
    r2 = r2_score(y_test, y_pred)

    # Definir el umbral de error aceptable (por ejemplo, 25%)
    threshold = 0.25
    errors = np.abs((y_test - y_pred) / y_test)
    correct_predictions = errors < threshold
    accuracy_within_threshold = np.mean(correct_predictions)
    print(f"Accuracy (dentro del umbral del {threshold * 100}%): {accuracy_within_threshold:.2f}")
    print(f"Mean Squared Error: {mse:.2f}")
    print(f"R^2 Score: {r2:.2f}")

    # Mostrar algunas predicciones
    comparison = pd.DataFrame({'Real': y_test, 'Predicted': y_pred})
    print(comparison.head(5))
    
