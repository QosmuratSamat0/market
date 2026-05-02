FROM golang:1.25-alpine AS builder

RUN apk add --no-cache --upgrade ca-certificates tzdata curl

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN mkdir -p /app/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/order-service ./cmd/main.go

FROM alpine:3.20

RUN apk add --no-cache --upgrade ca-certificates tzdata curl
WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/build/order-service ./order-service
COPY migrations ./migrations

RUN chown -R app:app /app
USER app

EXPOSE 8083

ENTRYPOINT ["./order-service"]
