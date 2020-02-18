#!/bin/bash

kubectl create -f deploy/charon.charon.cr_charons_crd.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/crds/charon.charon.cr_charons_crd.yaml
kubectl delete -f deploy/operator.yaml

kubectl create -f deploy/charon.charon.cr_charons_crd.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
kubectl create -f deploy/crds/charon.charon.cr_charons_crd.yaml
