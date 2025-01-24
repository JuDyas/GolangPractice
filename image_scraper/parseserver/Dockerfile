FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o parserserver main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/parserserver .

EXPOSE 8081

CMD ["./parserserver"]
