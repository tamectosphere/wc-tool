FROM golang:1.21-alpine as base

MAINTAINER "Pattadon Sa-ngasri<tam.ectosphere@gmail.com>"
WORKDIR /app

# Build pahse
FROM base as builder

RUN apk add --no-cache bash

COPY go.mod /app

RUN go mod tidy

COPY . /app

RUN go build cmd/wc/main.go

COPY docker-entrypoint.sh /app
RUN chmod +x /app/docker-entrypoint.sh

# ENTRYPOINT ["/app/docker-entrypoint.sh"]
# CMD ["default"]