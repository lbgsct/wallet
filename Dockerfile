FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wallet-service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wallet-service .

COPY config.env .

EXPOSE 8080

CMD ["./wallet-service"]
