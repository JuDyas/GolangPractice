FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o masterserver ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/masterserver .

EXPOSE 8080

CMD ["./masterserver"]
