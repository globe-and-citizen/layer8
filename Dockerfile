# Build the server
FROM golang:1.23.8-alpine3.20 AS builder

RUN mkdir /build

COPY ./server /build
COPY ./migrations /migrations

WORKDIR /build

RUN go get github.com/globe-and-citizen/layer8-utils

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go
RUN GOOS=linux GOARCH=amd64 go build -o setup cmd/setup/setup.go

# Deploy the server

FROM alpine:latest

RUN mkdir /layer8-app

WORKDIR /layer8-app

COPY --from=builder /migrations /layer8-app/migrations

COPY --from=builder /build/main /layer8-app

COPY --from=builder /build/assets-v1 /layer8-app/assets-v1

COPY --from=builder /build/certificates /layer8-app/certificates

COPY --from=builder /build/setup /layer8-app

RUN touch /layer8-app/.env

EXPOSE 5001

RUN chmod +x ./main
RUN chmod +x ./setup

# run commands will be specified in docker-compose.yml
#ENTRYPOINT ["./setup", "&&", "sleep 10", "&&", "./main"]