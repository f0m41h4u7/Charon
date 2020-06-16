FROM    golang:1.14.2-alpine3.11

ENV     GO111MODULE on

RUN	apk update && apk upgrade && apk add vim curl git make mercurial docker bash tree

RUN	echo && echo "Clone operator-sdk from github.com..." && echo && \
        mkdir -p $GOPATH/src/github.com/operator-framework && \
        cd $GOPATH/src/github.com/operator-framework && \
        git clone https://github.com/operator-framework/operator-sdk && \
        cd operator-sdk && \
        git fetch && git checkout && \
        echo && echo "make tidy..." && echo && \
        make tidy && \
        echo && echo "make install..." && echo && \
        make install

RUN	curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.15.0/bin/linux/amd64/kubectl \
	&& chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl

WORKDIR	/go/src/	

VOLUME	[ "/sys/fs/cgroup", "/go/src" ]
VOLUME	[ "operator-sdk", "/go/pkg"  ]
