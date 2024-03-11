# FROM golang:1.21.4-alpine3.17
FROM golang:1.21.6-alpine3.18

RUN mkdir /build

COPY ./server /build

WORKDIR /build

RUN go get github.com/globe-and-citizen/layer8-utils

RUN go mod tidy

RUN go build -o main .

EXPOSE 5001

RUN chmod +x ./main

ENTRYPOINT ["./main"]