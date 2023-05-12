# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /garybot .

EXPOSE 8080

ENTRYPOINT [ "/garybot" ]