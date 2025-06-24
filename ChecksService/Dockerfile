FROM golang:1.23.4 AS builder
WORKDIR /app

COPY . .
ARG CGO_ENABLED=0 \
    GOOS=linux

RUN go mod tidy && go build -o ./bin/app ./cmd/app.go

FROM alpine:latest
COPY --from=builder /app/bin/app /

EXPOSE 3000

ENTRYPOINT ["/app"]