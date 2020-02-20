#!/bin/bash

kubectl delete -f grafana-datasource-config.yaml
kubectl delete -f deployment-grafana.yaml
kubectl delete -f service-grafana.yaml

kubectl create -f grafana-datasource-config.yaml
kubectl create -f deployment-grafana.yaml
kubectl create -f service-grafana.yaml
