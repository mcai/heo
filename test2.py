import numpy as np
import pandas as pd
from keras.layers import LSTM, Dense
from keras.metrics import top_k_categorical_accuracy
from keras.models import Sequential
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, OneHotEncoder


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
model.add(LSTM(units=25, return_sequences=True, input_shape=(sequence_length, num_features)))
model.add(LSTM(units=25))
model.add(Dense(units=num_classes, activation='softmax'))
model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=[
    'accuracy', top_5_accuracy, top_10_accuracy
])
model.summary()

model.fit(train_X, train_Y, batch_size=1, epochs=5)

test_X = test[:, :, :-1]
test_Y = test[:, :, -1:]

test_Y = one_hot_encoder_address_delta.transform(test_Y.reshape(test_Y.shape[0], test_Y.shape[1]))

result = model.evaluate(test_X, test_Y, verbose=1)

print("loss: %s, acc: %s, top_5_accuracy: %s, top_10_accuracy: %s" %(result[0], result[1], result[2], result[3]))
