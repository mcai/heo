import pandas as pd
from keras import Sequential
from keras.layers import LSTM, Dense
from keras_preprocessing.sequence import TimeseriesGenerator
from matplotlib import pyplot
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler


def plot_data_frame(df):
    # df = df.head(100000)

    df.to_csv(
        file_name + "_modified.txt",
        index=False,
        # header=False,
        columns=["id", "thread_id", "pc", "type", "data_address", "data_address_delta"]
    )

    print(df)

    df.plot(x='id', y=['pc', 'data_address'], style='o')
    pyplot.show()

    # df.plot(x='id', y=['pc'], style='o')
    # pyplot.show()

    # df.plot(x='id', y=['data_address'], style='o')
    # pyplot.show()
    #
    df.plot(x='id', y=['data_address_delta'], style='o')
    pyplot.show()

    # TODO: to be imported into keras/pytorch LSTM model


if __name__ == "__main__":
    file_name = "test_results/real/mst_ht/l2_requests_trace.txt"

    df = pd.read_csv(file_name, names=["thread_id", "pc", "type", "data_address"])

    df['pc'] = df['pc'].apply(int, base=16).apply(lambda x: int(x / 8))
    df['data_address'] = df['data_address'].apply(int, base=16).apply(lambda x: int(x / 8))
    df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
    df['id'] = df.index

    df = df[df['data_address_delta'].notnull()]

    # plot_data_frame(df)

    scaler = MinMaxScaler(feature_range=(0, 1))
    df = scaler.fit_transform(df[['thread_id', 'pc', 'data_address_delta']])

    train, test = train_test_split(df, test_size=0.15)

    look_back = 10
    n_features = 2

    train_data_gen = TimeseriesGenerator(train, train, length=look_back, sampling_rate=1,stride=1, batch_size=3)
    test_data_gen = TimeseriesGenerator(test, test, length=look_back, sampling_rate=1,stride=1, batch_size=1)

    model = Sequential()
    model.add(LSTM(25, input_shape=(look_back, n_features)))
    model.add(Dense(n_features, activation='softmax'))
    model.compile(loss='categorical_crossentropy', optimizer='adam', metrics=['acc'])
    model.summary()

    model.fit_generator(train_data_gen, epochs=100, steps_per_epoch=10)

    model.evaluate_generator(test_data_gen)

    print()
