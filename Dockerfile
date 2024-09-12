# Build stage
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -v ./...

RUN go build -o main ./cmd

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]