FROM   		golang:1.13-alpine3.10

ENV     	GO111MODULE on

RUN     	apk update && apk upgrade && apk add vim bash && \
		go get -u github.com/gin-gonic/gin

#COPY		deployer.go $GOPATH
