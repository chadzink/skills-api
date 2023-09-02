FROM golang:1.21

WORKDIR /usr/src/app/api

COPY ./api .

RUN go install github.com/cosmtrek/air@latest

RUN go mod tidy