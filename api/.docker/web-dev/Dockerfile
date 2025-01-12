FROM golang:1.23 AS dev

RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

EXPOSE 8080

CMD ["air"]
