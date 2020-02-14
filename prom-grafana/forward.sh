#!/bin/bash
kubepfm --target $(kubectl get pods -n default | grep grafana | awk '{print $1}'):3001:3000
