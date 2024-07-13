from Model import *
import os.path

def MTLF(model):
    if os.path.isfile('model.h5'):
        model = load_model('model.h5')
    print("MTLF")
    print("xtrain: ", "dim: ", x_train.shape, x_train, "ytrain: ", x_test)
    model.fit(x_train, y_train, epochs=2)
    model.save('model.h5')
    print("trainig finish")
