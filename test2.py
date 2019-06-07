import numpy as np
import pandas as pd
from keras.layers import LSTM, Dense
from keras.models import Sequential
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder

sequence_length = 1
num_features = 2
num_outputs = 1

model = Sequential()
model.add(LSTM(units=20, return_sequences=True, input_shape=(sequence_length, num_features)))
model.add(LSTM(units=20))
model.add(Dense(num_outputs))
model.compile(loss='mae', optimizer='adam', metrics=['acc'])
model.summary()

file_name = "test_results/real/mst_ht/l2_requests_trace.txt"

df = pd.read_csv(file_name, names=["thread_id", "pc", "type", "data_address"])

df['pc'] = df['pc'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address'] = df['data_address'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
df['id'] = df.index

df = df[["thread_id", "pc", "data_address_delta"]]

df = df[df['data_address_delta'].notnull()]

# TODO: embedding, concatenate, predict 10 classes

# for i in range(1, sequence_length):
#     df['data_address_delta_prev_' + str(i)] = df['data_address_delta'].shift(-i)
#
# df = df[df['data_address_delta'].notnull()]
# for i in range(1, sequence_length):
#     df = df[df['data_address_delta_prev_' + str(i)].notnull()]

df = df.values

df = df.reshape(np.size(df, 0), 1, np.size(df, 1))

X = df[:, :-1]
Y = df[:, -1:]

train, test = train_test_split(df, test_size=0.15)

encoder_pc = LabelEncoder()
encoder_data_address_delta = LabelEncoder()

train[:, :, 1] = encoder_pc.fit_transform(train[:, :, 1]).reshape(-1, 1)
train[:, :, -1] = encoder_data_address_delta.fit_transform(train[:, :, -1]).reshape(-1, 1)

train_X = train[:, :, :-1]
train_Y = train[:, :, -1:]

train_Y = train_Y.reshape(np.size(train_Y, 0), np.size(train_Y, 1))

model.fit(train_X, train_Y, batch_size=1, epochs=3)

test[:, :, 1] = encoder_pc.fit_transform(test[:, :, 1]).reshape(-1, 1)
test[:, :, -1] = encoder_data_address_delta.fit_transform(test[:, :, -1]).reshape(-1, 1)

test_X = test[:, :, :-1]
test_Y = test[:, :, -1:]

test_Y = test_Y.reshape(np.size(test_Y, 0), np.size(test_Y, 1))

model.evaluate(test_X, test_Y, verbose=0)
