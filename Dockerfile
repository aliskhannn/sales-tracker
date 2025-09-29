FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sales-tracker ./cmd/sales-tracker/main.go

EXPOSE 8080