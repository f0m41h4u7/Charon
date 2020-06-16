SHELL := bash
D_IMG := charon-deployer
A_IMG := charon-analyzer
OP_IMG := charon-operator
CWD := $(shell pwd)

deploy:
	kubectl apply -f deployer-configmap.yaml
	kubectl apply -f charon-operator/deploy/service_account.yaml
	kubectl apply -f charon-operator/deploy/role.yaml
	kubectl apply -f charon-operator/deploy/role_binding.yaml
	kubectl apply -f charon-operator/deploy/operator.yaml
	kubectl apply -f charon-operator/deploy/crds/app.custom.cr_apps_crd.yaml
	kubectl apply -f charon-operator/deploy/crds/deployer.charon.cr_deployers_crd.yaml
	kubectl apply -f deployer/deployer.yaml

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd charon-deployer/
	golangci-lint run ./...
	cd ../deployer/
	golangci-lint run .
	cd ../analyzer/
	golangci-lint run .

build:
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/deployer/$(D_IMG)/g' deployer/deployer.yaml
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/analyzer/$(A_IMG)/g' deployer/deployer.yaml
	docker build -t $(A_IMG) analyzer/
	docker build -t $(D_IMG) deployer/
	sed -i 's/docker.pkg.github.com\/f0m41h4u7\/charon\/operator/$(OP_IMG)/g' charon-operator/deploy/operator.yaml
	docker build -t operator charon-operator/operator-sdk-docker/
	cd charon-operator/
	docker run --rm \
		-v $(CWD):/go/src \
		-v //var/run/docker.sock:/var/run/docker.sock \
		operator \
		/bin/bash -c "cd charon-operator && \
		operator-sdk generate k8s && \
		operator-sdk build $(OP_IMG) && \
		rm -rf build/_output"

.PHONY: deploy lint build
