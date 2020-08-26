FROM   		  golang:1.14.6
ARG         APP
WORKDIR     /home
ENV         GO111MODULE=on
COPY        go.mod .
COPY        go.sum .
RUN         go mod download
COPY        . .
RUN         go build ./cmd/$APP/main.go
ENTRYPOINT	["./main"]
