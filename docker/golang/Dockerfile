FROM golang:alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache bash \
    && apk add --no-cache make protobuf-dev

WORKDIR /backend

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go install -mod=mod github.com/golang/protobuf/protoc-gen-go

ENTRYPOINT CompileDaemon -command="./proxy-checker-server"