FROM   		  golang:1.16.5
ARG         APP
WORKDIR     /app
COPY        go.mod .
COPY        go.sum .
RUN         go mod download
COPY        . .
RUN         go build -o app ./cmd/$APP/main.go
ENTRYPOINT	["./app"]
