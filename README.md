# Charon

[![Go](https://img.shields.io/github/go-mod/go-version/f0m41h4u7/Charon?filename=deployer%2Fgo.mod)](https://github.com/f0m41h4u7/Charon/blob/master/deployer/go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/f0m41h4u7/Charon)](https://goreportcard.com/report/github.com/f0m41h4u7/Charon)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Smart version managing system for K8s

* Automatically deploys a pod from image pushed to docker registry
* Manages pod lifecycle
* Analyzes pod metrics
* If an anomaly was detected, rolls the pod back to a stable version

## Install

* Create a .env file providing Prometheus and Docker registry URLs (See [.env.example](.env.example))
* Run `setup.sh`

## Charon architecture

![alt text](https://raw.githubusercontent.com/f0m41h4u7/Charon/master/charon-project-scheme.png)
