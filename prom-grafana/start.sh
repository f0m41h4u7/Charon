#!/bin/bash

# prometheus
kubectl create -f clusterRole.yaml
kubectl create -f config-map.yaml
kubectl create -f prometheus-deployment.yaml

# grafana
kubectl create -f grafana-datasource-config.yaml
kubectl create -f deployment.grafana.yaml
kubectl create -f service-grafana.yaml
