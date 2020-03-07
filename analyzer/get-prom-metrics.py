from prometheus_api_client import Metric, MetricsList, PrometheusConnect
from prometheus_api_client.utils import parse_datetime, parse_timedelta
import pandas as pd
import numpy as np

def calc_delta(vals):
	diff = vals - np.roll(vals, 1)
	diff[0] = 0
	return diff

def monotonically_inc(vals):
	# check corner case
	if len(vals) == 1:
		return True
	diff = calc_delta(vals)
	diff[np.where(vals == 0)] = 0

	if ((diff < 0).sum() == 0):
		return True
	else:
		return False

def get_metrics():
    pc = PrometheusConnect(url="http://prometheus-service:9090", disable_ssl=True)

    start_time = parse_datetime("1h")
    end_time = parse_datetime("now")
    chunk_size = parse_timedelta("now", "1h")

    dataset = []

    all = pc.all_metrics()
    for metric_type in all:
        metric_data = pc.get_metric_range_data(
		metric_type,
		start_time=start_time,
		end_time=end_time,
		chunk_size=chunk_size,
	)
	metrics_object_list = MetricsList(metric_data)
	for item in metrics_object_list:
		vals = np.array(item.metric_values['y'].tolist())
		if monotonically_inc(vals):
			vals = calc_delta(vals)

		df = pd.DataFrame()
		df['ds'] = item.metric_values['ds']
		df['y'] = vals
        dataset.append(df)

    return dataset
