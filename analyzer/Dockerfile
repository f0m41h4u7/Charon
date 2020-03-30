FROM    	golang:1.13-alpine3.10
ENV     	GO111MODULE on
COPY		analyzer.go $GOPATH
ENTRYPOINT	["go", "run", "analyzer.go"]
