import numpy as np
import pandas as pd
from keras.layers import LSTM, Dense, Dropout
from keras.metrics import top_k_categorical_accuracy
from keras.models import Sequential
from keras.utils import plot_model
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, OneHotEncoder
import matplotlib.pyplot as plt


def top_5_accuracy(y_true, y_pred):
    return top_k_categorical_accuracy(y_true, y_pred, k=5)


def top_10_accuracy(y_true, y_pred):
    return top_k_categorical_accuracy(y_true, y_pred, k=10)


sequence_length = 1
num_features = 2

file_name = 'test_results/real/mst_ht/l2_requests_trace.txt'

df = pd.read_csv(file_name, names=['thread_id', 'pc', 'type', 'data_address'])

df['pc'] = df['pc'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address'] = df['data_address'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
df['id'] = df.index

df = df[['thread_id', 'pc', 'data_address_delta']]

df = df[df['data_address_delta'].notnull()]

# for i in range(1, sequence_length):
#     df['data_address_delta_prev_' + str(i)] = df['data_address_delta'].shift(-i)
#
# df = df[df['data_address_delta'].notnull()]
# for i in range(1, sequence_length):
#     df = df[df['data_address_delta_prev_' + str(i)].notnull()]

df = df.values

df = df.reshape(df.shape[0], 1, df.shape[1])

encoder_pc = LabelEncoder()
encoder_data_address_delta = LabelEncoder()

df[:, :, 1] = encoder_pc.fit_transform(df[:, :, 1]).reshape(-1, 1)
df[:, :, -1] = encoder_data_address_delta.fit_transform(df[:, :, -1]).reshape(-1, 1)

train, test = train_test_split(df, test_size=0.3)

train_X = train[:, :, :-1]
train_Y = train[:, :, -1:]

one_hot_encoder_address_delta = OneHotEncoder(handle_unknown='ignore')

train_Y = one_hot_encoder_address_delta.fit_transform(train_Y.reshape(train_Y.shape[0], train_Y.shape[1]))

num_classes = np.size(train_Y, -1)

model = Sequential()
model.add(LSTM(units=50, return_sequences=True, input_shape=(sequence_length, num_features)))
model.add(Dropout(0.2))
model.add(LSTM(units=50, return_sequences=True))
model.add(Dropout(0.2))
model.add(LSTM(units=50))
model.add(Dropout(0.2))
model.add(Dense(units=num_classes, activation='softmax'))
model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=[
    'accuracy', top_5_accuracy, top_10_accuracy
])
model.summary()

# plot_model(model, to_file='model.png')

history = model.fit(train_X, train_Y, batch_size=4, epochs=30, verbose=2, validation_split=0.25)

plt.plot(history.history['acc'])
plt.plot(history.history['val_acc'])
plt.title('Model Accuracy')
plt.ylabel('Accuracy')
plt.xlabel('Epoch')
plt.legend(['Train', 'Validation'], loc='upper left')
plt.show()

plt.plot(history.history['loss'])
plt.plot(history.history['val_loss'])
plt.title('Model Loss')
plt.ylabel('Loss')
plt.xlabel('Epoch')
plt.legend(['Train', 'Validation'], loc='upper left')
plt.show()

plt.plot(history.history['top_5_accuracy'])
plt.plot(history.history['val_top_5_accuracy'])
plt.title('Model Top 5 Accuracy')
plt.ylabel('Top 5 Accuracy')
plt.xlabel('Epoch')
plt.legend(['Train', 'Validation'], loc='upper left')
plt.show()

plt.plot(history.history['top_10_accuracy'])
plt.plot(history.history['val_top_10_accuracy'])
plt.title('Model Top 10 Accuracy')
plt.ylabel('Top 10 Accuracy')
plt.xlabel('Epoch')
plt.legend(['Train', 'Validation'], loc='upper left')
plt.show()

test_X = test[:, :, :-1]
test_Y = test[:, :, -1:]

test_Y = one_hot_encoder_address_delta.transform(test_Y.reshape(test_Y.shape[0], test_Y.shape[1]))

result = model.evaluate(test_X, test_Y, verbose=2)

print('loss: {:0.4f}, acc: {:0.4f}, top_5_accuracy: {:0.4f}, top_10_accuracy: {:0.4f}'.format(result[0], result[1], result[2], result[3]))
