#!/bin/bash

# prometheus
kubectl delete -f clusterRole.yaml
kubectl delete -f config-map.yaml
kubectl delete -f prometheus-deployment.yaml

kubectl create -f clusterRole.yaml
kubectl create -f config-map.yaml
kubectl create -f prometheus-deployment.yaml
