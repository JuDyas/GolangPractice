FROM golang:1.23.3
LABEL authors="judyas"

WORKDIR /pastebin
COPY . .
RUN go mod download
RUN go build -o /pastebin_app ./cmd/main.go
EXPOSE 8080
CMD ["/pastebin_app"]