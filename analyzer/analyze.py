#from forecast import ProphetForecast
from __future__ import absolute_import, division, print_function, unicode_literals
import tensorflow as tf
from get_prom_metrics import get_metrics

def create_time_steps(length):
    return list(range(-length, 0))

def show_plot(plot_data, delta, title):
    labels = ['History', 'True Future', 'Model Prediction']
    marker = ['.-', 'rx', 'go']
    time_steps = create_time_steps(plot_data[0].shape[0])
    if delta:
        future = delta
    else:
        future = 0

    plt.title(title)
    for i, x in enumerate(plot_data):
        if i:
            plt.plot(future, plot_data[i], marker[i], markersize=10, label=labels[i])
        else:
            plt.plot(time_steps, plot_data[i].flatten(), marker[i], label=labels[i])

    plt.legend()
    plt.xlim([time_steps[0], (future+5)*2])
    plt.xlabel('Time-Step')
    plt.savefig("plot.png")

df = get_metrics()
dataset = df.values
TRAIN_SPLIT = 70000

train_mean = dataset[:TRAIN_SPLIT].mean()
train_std = dataset[:TRAIN_SPLIT].std()
dataset = (dataset-train_mean)/train_std

past_history = 20
future_target = 72
EVALUATION_INTERVAL = 200
EPOCHS = 10
STEP = 6

x_train, y_train = univariate_data(data, 0, TRAIN_SPLIT, past_history, future_target)
x_val, y_val = univariate_data(data, TRAIN_SPLIT, None, past_history, future_target)

train_data_multi = tf.data.Dataset.from_tensor_slices((x_train, y_train))
train_data_multi = train_data_multi.cache().shuffle(BUFFER_SIZE).batch(BATCH_SIZE).repeat()

val_data_multi = tf.data.Dataset.from_tensor_slices((x_val, y_val))
val_data_multi = val_data.batch(BATCH_SIZE).repeat()

multi_step_model = tf.keras.models.Sequential()
multi_step_model.add(tf.keras.layers.LSTM(32, return_sequences=True, input_shape=x_train_multi.shape[-2:]))
multi_step_model.add(tf.keras.layers.LSTM(16, activation='relu'))
multi_step_model.add(tf.keras.layers.Dense(72))

multi_step_model.compile(optimizer=tf.keras.optimizers.RMSprop(clipvalue=1.0), loss='mae')

for x, y in val_data_multi.take(1):
    print (multi_step_model.predict(x).shape)

show_plot([x_train[0], y_train[0]], 0, 'Test metrics')
