import numpy as np
import pandas as pd
import tensorflow as tf
from keras import optimizers
from keras.layers import LSTM, Dense, Dropout
from keras.metrics import top_k_categorical_accuracy
from keras.models import Sequential
# from keras.utils import plot_model
from keras.utils import to_categorical
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, OneHotEncoder
# import matplotlib.pyplot as plt


class DeepNet:
    model = None
    batch_size = None
    epochs = None
    num_classes = None
    encoder_pc = None

    one_hot_encoder_data_address_delta = None

    def __init__(self, batch_size=1, epochs=1):
        self.batch_size = batch_size
        self.epochs = epochs

    def top_5_accuracy(self, y_true, y_pred):
        return top_k_categorical_accuracy(y_true, y_pred, k=5)

    def top_10_accuracy(self, y_true, y_pred):
        return top_k_categorical_accuracy(y_true, y_pred, k=10)

    def fit(self, file_name):
        sequence_length = 1
        num_features = 2

        df = pd.read_csv(file_name, names=['thread_id', 'pc', 'type', 'data_address'])

        df['pc'] = df['pc'].apply(int, base=16)
        df['data_address'] = df['data_address'].apply(int, base=16)
        df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
        # df['id'] = df.index

        df.to_csv(file_name + "_transformed2.txt")

        df = df[['thread_id', 'pc', 'data_address_delta']]

        df = df[df['data_address_delta'].notnull()]

        # for i in range(1, sequence_length):
        #     df['data_address_delta_prev_' + str(i)] = df['data_address_delta'].shift(-i)
        #
        # df = df[df['data_address_delta'].notnull()]
        # for i in range(1, sequence_length):
        #     df = df[df['data_address_delta_prev_' + str(i)].notnull()]

        df.to_csv(file_name + "_transformed.txt")

        df = df.values

        df = df.reshape(df.shape[0], 1, df.shape[1])

        self.encoder_pc = LabelEncoder() # TODO: should be used in predict too

        df[:, :, 1] = self.encoder_pc.fit_transform(df[:, :, 1]).reshape(-1, 1)

        train, test = train_test_split(df, test_size=0.3)

        train_X = train[:, :, :-1]
        train_Y = train[:, :, -1:]

        self.one_hot_encoder_data_address_delta = OneHotEncoder(handle_unknown='ignore')

        train_Y = self.one_hot_encoder_data_address_delta.fit_transform(train_Y.reshape(train_Y.shape[0], train_Y.shape[1]))

        self.num_classes = np.size(train_Y, -1)

        learning_rate = 0.001
        decay_rate = learning_rate / self.epochs

        self.model = Sequential()
        self.model.add(LSTM(units=128, return_sequences=True, input_shape=(sequence_length, num_features)))
        # self.model.add(Dropout(0.2))
        self.model.add(LSTM(units=128))
        # self.model.add(Dropout(0.2))
        self.model.add(Dense(units=self.num_classes, activation='softmax'))
        self.model.compile(optimizer=optimizers.Adam(lr=learning_rate, decay=decay_rate), loss='categorical_crossentropy', metrics=[
            'accuracy', self.top_5_accuracy, self.top_10_accuracy
        ])
        self.model.summary()

        # plot_model(model, to_file='model.png')

        history = self.model.fit(train_X, train_Y, batch_size=self.batch_size, epochs=self.epochs, verbose=2, validation_split=0.25)

        # plt.plot(history.history['acc'])
        # plt.plot(history.history['val_acc'])
        # plt.title('Model Accuracy')
        # plt.ylabel('Accuracy')
        # plt.xlabel('Epoch')
        # plt.legend(['Train', 'Validation'], loc='upper left')
        # plt.show()
        #
        # plt.plot(history.history['loss'])
        # plt.plot(history.history['val_loss'])
        # plt.title('Model Loss')
        # plt.ylabel('Loss')
        # plt.xlabel('Epoch')
        # plt.legend(['Train', 'Validation'], loc='upper left')
        # plt.show()
        #
        # plt.plot(history.history['top_5_accuracy'])
        # plt.plot(history.history['val_top_5_accuracy'])
        # plt.title('Model Top 5 Accuracy')
        # plt.ylabel('Top 5 Accuracy')
        # plt.xlabel('Epoch')
        # plt.legend(['Train', 'Validation'], loc='upper left')
        # plt.show()
        #
        # plt.plot(history.history['top_10_accuracy'])
        # plt.plot(history.history['val_top_10_accuracy'])
        # plt.title('Model Top 10 Accuracy')
        # plt.ylabel('Top 10 Accuracy')
        # plt.xlabel('Epoch')
        # plt.legend(['Train', 'Validation'], loc='upper left')
        # plt.show()

        test_X = test[:, :, :-1]
        test_Y = test[:, :, -1:]

        test_Y = self.one_hot_encoder_data_address_delta.transform(test_Y.reshape(test_Y.shape[0], test_Y.shape[1]))

        result = self.model.evaluate(test_X, test_Y, verbose=2)

        print('loss: {:0.4f}, acc: {:0.4f}, top_5_accuracy: {:0.4f}, top_10_accuracy: {:0.4f}'.format(result[0], result[1], result[2], result[3]))

    def predict(self, thread_id, pc, top_k=10):
        X = np.array([thread_id, pc])
        X = X.reshape(1, 1, X.shape[0]) # num_samples=1, num_steps=1, num_features=2

        X[:, :, 1] = self.encoder_pc.transform(X[:, :, 1]).reshape(-1, 1)

        predicted_Y = self.model.predict(X, batch_size=1, verbose=2)

        with tf.Session().as_default():
            scores, indices = tf.math.top_k(
                tf.convert_to_tensor(predicted_Y, dtype=np.float32),
                k=top_k,
                sorted=True,
            )

            scores = scores.eval()
            indices = indices.eval()

            xxx = to_categorical(indices, num_classes=self.num_classes)

            predictions = []

            for i in range(np.size(xxx, 1)):
                xx = xxx[:, i, :]
                inversed_xx = self.one_hot_encoder_data_address_delta.inverse_transform(xx)
                predictions.append(int(inversed_xx[0, 0]))

            return predictions


if __name__ == "__main__":
    deep_net = DeepNet(epochs=20)
    deep_net.fit("/Users/itecgo/go/src/github.com/mcai/heo/test_results/real/mst_ht/l2_requests_trace.txt")

    def predict(pc):
        predictions = deep_net.predict(0, pc)

        print("pc: ")
        print(pc)

        for prediction in predictions:
            print(prediction)

    predict(4196900)

    print()
    print()

    predict(4196916)

# TODO: @tragu in my case: reduced the number of output classes, increased the number of samples of each class and fixed inconsistencies using a pre-trained model
