FROM golang:1.24.2-alpine AS builder

WORKDIR /services

ARG CGO_ENABLED=0 \
    GOOS=linux

COPY . .
RUN go work vendor

RUN go build -o ./bin/checks -mod vendor ./ChecksService/cmd/app.go
RUN go build -o ./bin/warns -mod vendor ./WarnsService/cmd/app.go
RUN go build -o ./bin/promos -mod vendor ./PromosService/cmd/app.go

FROM alpine:latest
COPY --from=builder /services/bin/checks /services/bin/promos /services/bin/warns /