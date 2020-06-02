<p align="center">
	<img src="https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon.png"><br>
</p>
<h3 align="center">Smart version managing system for K8s</h3>
<p align="center">
	<a href="https://github.com/f0m41h4u7/Charon/blob/master/main/go.mod" rel="nofollow"><img alt="Go" src="https://img.shields.io/github/go-mod/go-version/f0m41h4u7/Charon?filename=deployer%2Fgo.mod"></a>
	<a href="https://goreportcard.com/report/github.com/f0m41h4u7/Charon" rel="nofollow"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/f0m41h4u7/Charon"></a>
	<a href="https://github.com/f0m41h4u7/Charon/LICENSE" rel="nofollow"><img alt="MIT License" src="https://img.shields.io/github/license/f0m41h4u7/Charon"></a>
</p>

## How it works

* Automatically deploys a pod from image pushed to docker registry
* Manages pod lifecycle
* Analyzes pod metrics
* If an anomaly was detected, rolls the pod back to a stable version

### Charon architecture

![alt text](https://raw.githubusercontent.com/f0m41h4u7/Charon/master/files/charon-project-scheme.png)

## Usage

* Provide Prometheus URL and Docker Registry name in `deployer-configmap.yaml`
* Configure your Docker Registry notification endpoint to send alerts to `http://your-url:31337/rollout`
* Run `make` to deploy Charon
* Run `make build` to build Charon from source
