FROM   		golang:1.13-alpine3.10
ENV     	GO111MODULE on
COPY		analyzer.go /home
COPY		promQuery.go /home
COPY		.env /home
COPY		go.mod /home
COPY		go.sum /home
WORKDIR		/home
ENTRYPOINT	["go", "run", "analyzer.go"]
