#!/bin/sh
operator-sdk generate k8s
operator-sdk build charon-operator
rm -rf charon-operator/build/_output
