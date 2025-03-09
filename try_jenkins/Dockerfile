FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /app .

EXPOSE 8081

CMD ["./main"]