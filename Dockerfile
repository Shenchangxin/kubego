FROM golang:alpine as builder
WORKDIR /go/src/kubego.com/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER = "282578874@qq.com"

WORKDIR /go/src/kubego.com/server
COPY --from=0 /go/src/kubego.com/server/config.yaml ./config.yaml
COPY --from=0 /go/src/kubego.com/server/.kube/config ./.kube/config
COPY --from=0 /go/src/kubego.com/server/server ./
EXPOSE 8082
ENTRYPOINT ./server