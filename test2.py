import pandas as pd
from keras.layers import LSTM, Dense
from keras.models import Sequential
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder

sequence_length = 1
num_outputs = 1

model = Sequential()
model.add(LSTM(20, return_sequences=True, input_shape=(sequence_length, 1)))
model.add(LSTM(20))
model.add(Dense(num_outputs))
model.compile(loss='mae', optimizer='adam', metrics=['acc'])
model.summary()

file_name = "test_results/real/mst_ht/l2_requests_trace.txt"

df = pd.read_csv(file_name, names=["thread_id", "pc", "type", "data_address"])

df['pc'] = df['pc'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address'] = df['data_address'].apply(int, base=16).apply(lambda x: int(x / 8))
df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
df['id'] = df.index

df = df[["pc", "data_address_delta"]]

df = df[df['data_address_delta'].notnull()]

# df = pd.concat([df,pd.get_dummies(df['pc'], prefix='pc')],axis=1)
# df.drop(['pc'],axis=1, inplace=True)
#
# df = pd.concat([df,pd.get_dummies(df['data_address_delta'], prefix='data_address_delta')],axis=1)
# df.drop(['data_address_delta'],axis=1, inplace=True)

# TODO: embedding, concatenate, predict 10 classes

# for i in range(1, sequence_length):
#     df['data_address_delta_prev_' + str(i)] = df['data_address_delta'].shift(-i)
#
# df = df[df['data_address_delta'].notnull()]
# for i in range(1, sequence_length):
#     df = df[df['data_address_delta_prev_' + str(i)].notnull()]

df = df.values

X = df[:, :-1]
Y = df[:, -1:]

train, test = train_test_split(df, test_size=0.15)

train_X = train[:, :-1]
train_Y = train[:, -1:]

encoder_X = LabelEncoder()
encoder_Y = LabelEncoder()

encoded_X = encoder_X.fit_transform(train_X).reshape(len(train_X), 1, 1)
encoded_Y = encoder_Y.fit_transform(train_Y).reshape(len(train_Y), 1)

model.fit(encoded_X, encoded_Y, batch_size=1, epochs=10)

test_X = test[:, :-1]
test_Y = test[:, -1:]

encoded_X = encoder_X.fit_transform(test_X).reshape(len(test_X), 1, 1)
encoded_Y = encoder_Y.fit_transform(test_Y).reshape(len(test_Y), 1)

model.evaluate(encoded_X, encoded_Y, verbose=0)
