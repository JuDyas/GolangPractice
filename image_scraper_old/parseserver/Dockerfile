FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o parseserver ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /app/parseserver .
COPY --from=builder /app/config ./config

ENV PARSER_SERVER_PORT=${PARSER_SERVER_PORT}
EXPOSE ${PARSER_SERVER_PORT}

CMD ["./parseserver"]
