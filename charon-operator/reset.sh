#!/bin/bash

kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/crds/charon.charon.cr_charons_crd.yaml
kubectl delete -f deploy/operator.yaml
