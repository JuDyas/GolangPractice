FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o parseserver ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/parseserver .

EXPOSE 8081

CMD ["./parseserver"]
