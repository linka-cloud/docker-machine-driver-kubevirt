FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /go/src/github.com/docker/machine

RUN git clone --branch=alpine-provisioning https://gitlab.com/linka-cloud/docker-machine . && \
    GO111MODULE=off go build -o /usr/local/bin/docker-machine ./cmd/docker-machine

WORKDIR /docker-machine-driver-kubevirt

ADD go.mod go.mod
ADD go.sum go.sum

RUN go mod download

ADD driver driver
ADD main.go main.go

RUN go build -o /usr/local/bin/docker-machine-driver-kubevirt .

FROM alpine

RUN apk add --no-cache openssh-client docker-cli curl

COPY --from=0 /usr/local/bin/docker-machine /usr/local/bin/docker-machine
COPY --from=0 /usr/local/bin/docker-machine-driver-kubevirt /usr/local/bin/docker-machine-driver-kubevirt

ENV SHELL=sh

CMD ["tail", "-f", "/dev/null"]
