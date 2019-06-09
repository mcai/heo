import numpy as np
import pandas as pd
import tensorflow as tf
from keras.layers import LSTM, Dense, Dropout
from keras.metrics import top_k_categorical_accuracy
from keras.models import Sequential
# from keras.utils import plot_model
from keras.utils import to_categorical
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, OneHotEncoder
# import matplotlib.pyplot as plt

encoder = OneHotEncoder(handle_unknown='ignore', sparse=False)

Y = np.array([[32.], [64.], [128.]])

transformed_Y = encoder.fit_transform(Y)

inversed_Y = encoder.inverse_transform(transformed_Y)

print("Y:")
print(Y)

print("transformed_Y:")
print(transformed_Y)

print("inversed_Y:")
print(inversed_Y)

predicted_Y = np.array([[0.3, 0.2, 0.5]])
print("predicted_Y:")
print(predicted_Y)

inversed_predicted_Y = encoder.inverse_transform(predicted_Y)

print("inversed_predicted_Y:")
print(inversed_predicted_Y)




sess = tf.Session()

with sess.as_default():
    scores, indices = tf.math.top_k(
        tf.convert_to_tensor(predicted_Y, dtype=np.float32),
        k=2,
        sorted=True,
    )

    scores = scores.eval()
    indices = indices.eval()

    xxx = to_categorical(indices)

    for i in range(np.size(xxx, 1)):
        xx = xxx[:, i, :]
        inversed_xx = encoder.inverse_transform(xx)
        print("inversed_xx:")
        print(inversed_xx)


print()