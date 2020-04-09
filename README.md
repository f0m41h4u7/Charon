# Charon

<p align="center">
	<img src="https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon.jpg"><br>
	<b>Smart version managing system for K8s</b><br>
	<a href="https://github.com/f0m41h4u7/Charon/blob/master/deployer/go.mod" rel="nofollow"><img alt="Go" src="https://img.shields.io/github/go-mod/go-version/f0m41h4u7/Charon?filename=deployer%2Fgo.mod">
	<a href="https://goreportcard.com/report/github.com/f0m41h4u7/Charon" rel="nofollow"><ing alt="Go Report Card" src="https://goreportcard.com/badge/github.com/f0m41h4u7/Charon">
	<a href="https://www.gnu.org/licenses/gpl-3.0"><img alt="License: GPL v3" src="https://img.shields.io/badge/License-GPLv3-blue.svg">
</p>

* Automatically deploys a pod from image pushed to docker registry
* Manages pod lifecycle
* Analyzes pod metrics
* If an anomaly was detected, rolls the pod back to a stable version

## Install

* Create a .env file providing Prometheus and Docker registry URLs (See [.env.example](.env.example))
* Run `setup.sh`

## Charon architecture

![alt text](https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon-project-scheme.png)
