FROM   		golang:1.13-alpine3.10
ENV     	GO111MODULE on
COPY		deployer.go $GOPATH
COPY		../.env $GOPATH
ENTRYPOINT	["go", "run", "deployer.go"]
