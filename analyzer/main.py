from forecast import ProphetForecast
from get_prom_metrics import get_metrics

df = get_metrics()

vals = df["y"]
train = df[0 : int(0.7 * len(vals))]
test = df[int(0.7 * len(vals)) :]

pf = ProphetForecast(train, test)
forecast = pf.fit_model(len(test))
pf.graph()
