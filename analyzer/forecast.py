from fbprophet import Prophet
import pandas as pd
import numpy as np
import matplotlib.pylab as plt
import datetime as dt


class ProphetForecast:
    def __init__(self, train, test):
        self.train = train
        self.test = test

    def fit_model(self, n_predict):
        m = Prophet(
            daily_seasonality=False, weekly_seasonality=False, yearly_seasonality=False
        )
        m.fit(self.train)
        future = m.make_future_dataframe(periods=len(self.test), freq="1MIN")
        self.forecast = m.predict(future)

        return self.forecast

    def graph(self):
        fig = plt.figure(figsize=(40, 10))
        plt.plot(
            np.array(self.train["ds"]),
            np.array(self.train["y"]),
            "b",
            label="train",
            linewidth=3,
        )
        plt.plot(
            np.array(self.test["ds"]),
            np.array(self.test["y"]),
            "g",
            label="test",
            linewidth=3,
        )

        forecast_ds = np.array(self.forecast["ds"])
        plt.plot(
            forecast_ds, np.array(self.forecast["yhat"]), "o", label="yhat", linewidth=3
        )
        plt.plot(
            forecast_ds,
            np.array(self.forecast["yhat_upper"]),
            "y",
            label="yhat_upper",
            linewidth=3,
        )
        plt.plot(
            forecast_ds,
            np.array(self.forecast["yhat_lower"]),
            "y",
            label="yhat_lower",
            linewidth=3,
        )
        plt.xlabel("Timestamp")
        plt.ylabel("Value")
        plt.legend(loc=1)
        plt.title("Prophet Model Forecast")
        plt.savefig("plot.png")
