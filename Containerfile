FROM docker.io/library/golang:alpine as builder
WORKDIR /build
ADD . /build
RUN go version && \
    go env && \
    CGO_ENABLED=0 GOOS=linux go build

FROM docker.io/library/alpine:3
MAINTAINER sethops1.net

WORKDIR /opt
COPY --from=builder /build/nomad-holepunch /opt

ENTRYPOINT ["/opt/nomad-holepunch"]
