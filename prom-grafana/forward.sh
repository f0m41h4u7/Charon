#!/bin/bash
kubepfm --target $(kubectl get pods -n default | grep prom | awk '{print $1}'):8080:9090 --target $(kubectl get pods -n default | grep grafana | awk '{print $1}'):3001:3000
