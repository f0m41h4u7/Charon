#!/bin/sh
kubectl apply -f /opt/charon/config.yaml
kubectl apply -f /opt/charon/role.yaml
kubectl apply -f /opt/charon/role_binding.yaml
kubectl apply -f /opt/charon/service_account.yaml
kubectl apply -f /opt/charon/operator.yaml
kubectl apply -f /opt/charon/crds/
kubectl apply -f /opt/charon/deployer-analyzer.yaml
