FROM golang:1.12
MAINTAINER slimm609

ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 8080

RUN apt-get update && \
    apt-get install -y --no-install-suggests --no-install-recommends jq

WORKDIR /summon-ssm-secrets

RUN go get -u github.com/jstemmer/go-junit-report && \
    go get github.com/smartystreets/goconvey && \
    mkdir -p /summon-ssm-secrets/output

COPY go.mod go.sum /summon-ssm-secrets/
RUN go mod download

COPY . .
