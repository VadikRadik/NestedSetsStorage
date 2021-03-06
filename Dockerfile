FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /storage/build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build storage.go

EXPOSE 7090
