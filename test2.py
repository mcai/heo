import numpy as np
import pandas as pd
import tensorflow as tf
from keras.layers import LSTM, Dense, Dropout
from keras.metrics import top_k_categorical_accuracy
from keras.models import Sequential
# from keras.utils import plot_model
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, OneHotEncoder
# import matplotlib.pyplot as plt


class DeepNet:
    model = None
    batch_size = None
    epochs = None
    num_classes = None

    one_hot_encoder_address_delta = None

    def __init__(self, batch_size=4, epochs=30):
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

        self.one_hot_encoder_address_delta = OneHotEncoder(handle_unknown='ignore')

        train_Y = self.one_hot_encoder_address_delta.fit_transform(train_Y.reshape(train_Y.shape[0], train_Y.shape[1]))

        self.num_classes = np.size(train_Y, -1)

        self.model = Sequential()
        self.model.add(LSTM(units=50, return_sequences=True, input_shape=(sequence_length, num_features)))
        self.model.add(Dropout(0.2))
        self.model.add(LSTM(units=50, return_sequences=True))
        self.model.add(Dropout(0.2))
        self.model.add(LSTM(units=50))
        self.model.add(Dropout(0.2))
        self.model.add(Dense(units=self.num_classes, activation='softmax'))
        self.model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=[
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

        test_Y = self.one_hot_encoder_address_delta.transform(test_Y.reshape(test_Y.shape[0], test_Y.shape[1]))

        result = self.model.evaluate(test_X, test_Y, verbose=2)

        print('loss: {:0.4f}, acc: {:0.4f}, top_5_accuracy: {:0.4f}, top_10_accuracy: {:0.4f}'.format(result[0], result[1], result[2], result[3]))

    def predict(self, thread_id, pc, top_k=10):
        X = np.array([thread_id, pc])
        X = X.reshape(1, 1, X.shape[0]) # num_samples=1, num_steps=1, num_features=2

        # predictions = self.model.predict_classes(X, batch_size=1, verbose=2)
        # prediction_ = np.argmax(to_categorical(predictions), axis = 1)
        # prediction_ = encoder.inverse_transform(prediction_)
        #
        # for i, j in zip(prediction_ , predict_species):
        #     print( " the nn predict {}, and the species to find is {}".format(i,j))

        predicted_Y = self.model.predict(X, batch_size=1, verbose=2)

        print("predicted_Y: ")
        print(predicted_Y.shape)
        print(predicted_Y)

        sess = tf.Session()

        with sess.as_default():
            scores, indices = tf.math.top_k(
                tf.convert_to_tensor(predicted_Y, dtype=np.float32),
                k=top_k,
                sorted=True,
            )

            scores = scores.eval()
            indices = indices.eval()

            print("scores.shape")
            print(scores.shape)
            print(scores)

            print("indices.shape")
            print(indices.shape)
            print(indices)

            restored_predicted_Y = self.one_hot_encoder_address_delta.inverse_transform(predicted_Y)

            print("restored_predicted_Y.shape")
            print(restored_predicted_Y.shape)
            print(restored_predicted_Y)

            return restored_predicted_Y[0]
