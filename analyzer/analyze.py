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
print(df)

vals = df["y"]
train = df[0 : int(0.7 * len(vals))]
test = df[int(0.7 * len(vals)) :]

train_mean = train.mean()
train_std = train.std()
data = (train['y']-train_mean)/train_std

past_history = 20
future_target = 0

x_train, y_train = univariate_data(data, 0, TRAIN_SPLIT, past_history, future_target)
x_val, y_val = univariate_data(data, TRAIN_SPLIT, None, past_history, future_target)

print ('Single window of past history')
print (x_train[0])
print ('\n Target temperature to predict')
print (y_train[0])

show_plot([x_train[0], y_train[0]], 0, 'Test metrics')

#pf = ProphetForecast(train, test)
#forecast = pf.fit_model(len(test))
#pf.graph()
