from prometheus_api_client import Metric, MetricsList, PrometheusConnect
from prometheus_api_client.utils import parse_datetime, parse_timedelta

import matplotlib.pyplot as plt

pc = PrometheusConnect(url="http://167.172.137.177:30329", disable_ssl=True)

start_time = parse_datetime("3d")
end_time = parse_datetime("now")
chunk_size = parse_timedelta("now", "1d")

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
		print(item.metric_name, item.label_config)
		print(item.metric_values)
	print("\n ******************** \n")
