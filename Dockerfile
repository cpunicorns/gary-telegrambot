# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum *.go ./

RUN go mod download && GOOS=linux go build -o /garybot .

EXPOSE 8080

ENTRYPOINT [ "/garybot" ]