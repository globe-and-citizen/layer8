# Build the server
FROM golang:1.21.6-alpine3.18 AS builder

RUN mkdir /build

COPY ./server /build

WORKDIR /build

RUN go get github.com/globe-and-citizen/layer8-utils

RUN go mod tidy

RUN go build -o main .

# Deploy the server

FROM alpine:latest

RUN mkdir /layer8-app

WORKDIR /layer8-app

COPY --from=builder /build/main /layer8-app

COPY --from=builder /build/assets-v1 /layer8-app/assets-v1

COPY --from=builder /build/certificates /layer8-app/certificates

COPY --from=builder /build/.env /layer8-app

EXPOSE 5001

RUN chmod +x ./main

ENTRYPOINT ["./main"]