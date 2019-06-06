import pandas as pd
import numpy as np
from matplotlib import pyplot

if __name__ == "__main__":
    file_name = "test_results/real/mst_ht/l2_requests_trace.txt"

    df = pd.read_csv(file_name, names=["thread_id", "pc", "type", "data_address"])

    df['pc'] = df['pc'].apply(int, base=16).apply(lambda x: int(x / 8))
    df['data_address'] = df['data_address'].apply(int, base=16).apply(lambda x: int(x / 8))
    df['data_address_delta'] = df.groupby('pc', sort=False)['data_address'].diff()
    df['id'] = df.index

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
