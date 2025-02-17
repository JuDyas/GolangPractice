FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o masterserver ./cmd/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/masterserver .
COPY --from=builder /app/uploads ${UPLOADS_DIR}

ENV MASTER_SERVER_PORT=${MASTER_SERVER_PORT}
EXPOSE ${MASTER_SERVER_PORT}

CMD ["./masterserver"]
