SHELL := bash
D_IMG := charon-deployer
A_IMG := charon-analyzer
OP_IMG := charon-operator
CWD := $(shell pwd)

deploy:
	kubectl apply -f deployer-configmap.yaml
	kubectl apply -f deploy/

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

build:
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/deployer/$(D_IMG)/g' deploy/deployer.yaml
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/analyzer/$(A_IMG)/g' deploy/deployer.yaml
	docker build -t $(A_IMG) build/analyzer/
	docker build -t $(D_IMG) build/deployer/
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/operator/$(OP_IMG)/g' operator/deploy/operator.yaml
	docker build -t operator operator/operator-sdk-docker/
	cd operator/
	docker run --rm \
		-v $(CWD):/go/src \
		-v //var/run/docker.sock:/var/run/docker.sock \
		operator \
		/bin/bash -c "cd operator && \
		operator-sdk generate k8s && \
		operator-sdk build $(OP_IMG) && \
		rm -rf build/_output"

.PHONY: deploy lint build
