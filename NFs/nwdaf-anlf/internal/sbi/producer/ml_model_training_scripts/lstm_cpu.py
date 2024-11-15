import numpy as np
import pandas as pd
from sklearn.preprocessing import MinMaxScaler, StandardScaler
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
import matplotlib.pyplot as plt
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dense, Dropout
from tensorflow.keras.callbacks import EarlyStopping
from sklearn.model_selection import train_test_split
from tensorflow.keras.preprocessing.sequence import TimeseriesGenerator


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

# Redimensionar X para que sea compatible con la entrada del modelo LSTM (samples, time_steps, features)
# X = X.reshape((X.shape[0], X.shape[1], 1))

# Paso 3: Definir y entrenar el modelo LSTM
## Configuration 1
# model = Sequential()
# model.add(LSTM(units=100, return_sequences=False, input_shape=(time_steps, 1)))
# model.add(Dense(units=1))  # Capa de salida para predecir un valor continuo (uso de CPU)
# model.compile(optimizer='adam', loss='mse')

## Configuration 2
# Definir la arquitectura del modelo
model = Sequential()
# Primera capa LSTM
model.add(LSTM(50, activation='tanh', input_shape=(3, 1), return_sequences=True))
model.add(Dropout(0.2))  # Regularización para evitar sobreajuste
# Segunda capa LSTM
model.add(LSTM(50, activation='tanh', return_sequences=False))
model.add(Dropout(0.2))
# Capa de salida
model.add(Dense(1))  # Predicción del valor de uso de CPU en el siguiente paso
# Compilar el modelo
model.compile(optimizer='adam', loss='mse')

early_stopping = EarlyStopping(
    monitor='val_loss',      # Métrica que se va a monitorear
    patience=5,              # Esperar 10 epochs sin mejora
    restore_best_weights=True,  # Restaurar el mejor modelo al final
    mode='min',              # Minimizar la métrica (val_loss)
    verbose=1                # Mostrar mensajes en la consola
)

# Entrenar el modelo
history = model.fit(X_train,
                    y_train, 
                    epochs=30, 
                    batch_size=1,
                    validation_data=(X_test, y_test), 
                    verbose=1,
                    # validation_split=0.3,       # Utilizar el 20% de los datos para validación
                    callbacks=[early_stopping]  # Pasar el callback de Early Stopping
                    )

loss = model.evaluate(X_test, y_test)
print("MSE en conjunto de prueba:", loss)

# Predicción de la serie temporal
predictions = model.predict(X_train)

# Paso 4: Hacer predicciones con el modelo entrenado
# predictions = scaler.inverse_transform(predictions)

# Mostrar los valores reales y predichos
real_values = scalerMinMax.inverse_transform(y_train)

# Calcular las métricas
mse = mean_squared_error(real_values, predictions)
mae = mean_absolute_error(real_values, predictions)
r2 = r2_score(real_values, predictions)

# Mostrar las métricas
print(f"MSE: {mse:.4f}")
print(f"MAE: {mae:.4f}")
print(f"R² Score: {r2:.4f}")

# print("Valores reales:", real_values.flatten())
# print("Valores predichos:", predictions.flatten())

# Graficar los resultados
plt.figure(figsize=(10, 6))
plt.plot(real_values, label='Valores Reales')
plt.plot(predictions, label='Valores Predichos')
plt.title('Predicción del uso de CPU con LSTM')
plt.xlabel('Tiempos')
plt.ylabel('Uso de CPU')
plt.legend()
plt.grid(True)
plt.show()
