FROM golang:1.25.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

EXPOSE 8084

CMD ["./main"]
