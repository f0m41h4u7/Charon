<p align="center">
	<img src="https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon.png"><br>
</p>
<h3 align="center">Smart version managing system for K8s</h3>
<p align="center">
	<a href="https://github.com/f0m41h4u7/Charon/blob/master/main/go.mod" rel="nofollow"><img alt="Go" src="https://img.shields.io/github/go-mod/go-version/f0m41h4u7/Charon?filename=main%2Fgo.mod">
	<a href="https://goreportcard.com/report/github.com/f0m41h4u7/Charon" rel="nofollow"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/f0m41h4u7/Charon">
	<img alt="License: GPL v3" src="https://img.shields.io/github/license/f0m41h4u7/Charon">
</p>

## How does it work

* Automatically deploys a pod from image pushed to docker registry
* Manages pod lifecycle
* Analyzes pod metrics
* If an anomaly was detected, rolls the pod back to a stable version

### Charon architecture

![alt text](https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon-project-scheme.png)

## How to use

* Create a .env file providing Prometheus and Docker registry URLs (See [.env.example](.env.example))
* Run `setup.sh`
