SHELL := bash
D_IMG := charon-deployer
A_IMG := charon-analyzer
OP_IMG := charon-operator
CWD := $(shell pwd)

deploy:
	kubectl apply -f deployer-configmap.yaml
	kubectl apply -f operator/deploy/
	kubectl apply -f deploy/

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

build:
	go build -o $(A_IMG) ./cmd/analyzer/main.go
	go build -o $(D_IMG) ./cmd/deployer/main.go

build-docker:
	make analyzer
	make deployer

analyzer:
	docker build -t $(A_IMG) --build-arg APP=analyzer -f build/Dockerfile .

deployer:
	docker build -t $(D_IMG) --build-arg APP=deployer -f build/Dockerfile .

.PHONY: deploy lint build analyzer deployer operator
