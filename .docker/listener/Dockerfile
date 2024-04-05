# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS base

COPY ./listener /go/src
WORKDIR /go/src

RUN go mod download

#SET PATH ENV
ENV PATH="$PATH:/go/src"

RUN CGO_ENABLED=0 GOOS=linux go build -o /listener

CMD ["/listener"]
