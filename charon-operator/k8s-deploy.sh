#!/bin/bash

kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
kubectl create -f deploy/crds/charon.charon.cr_charons_crd.yaml
