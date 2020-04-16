#!/bin/bash

cp .env main/
cp .env analyzer/

kubectl apply -f charon-operator/deploy/service_account.yaml
kubectl apply -f charon-operator/deploy/role.yaml
kubectl apply -f charon-operator/deploy/role_binding.yaml
kubectl apply -f charon-operator/deploy/operator.yaml
kubectl apply -f charon-operator/deploy/crds/app.custom.cr_apps_crd.yaml
kubectl apply -f charon-operator/deploy/crds/deployer.charon.cr_deployers_crd.yaml
kubectl apply -f deployer/deployer.yaml
