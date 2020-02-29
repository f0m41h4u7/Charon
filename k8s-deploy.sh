#!/bin/bash

kubectl create -f charon-operator/deploy/service_account.yaml
kubectl create -f charon-operator/deploy/role.yaml
kubectl create -f charon-operator/deploy/role_binding.yaml
kubectl create -f charon-operator/deploy/operator.yaml
kubectl create -f charon-operator/deploy/crds/app.custom.cr_apps_crd.yaml
kubectl create -f charon-operator/deploy/crds/deployer.charon.cr_deployers_crd.yaml
kubectl create -f deployer/cr.yaml
