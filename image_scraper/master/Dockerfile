FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o master ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/master .
COPY --from=builder /app/internal/db/migrations /migrations


EXPOSE ${MASTER_SERVER_PORT}

CMD ["./master"]