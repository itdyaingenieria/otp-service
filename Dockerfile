# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o otp-service ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/otp-service .
COPY migrations ./migrations
COPY .env .
EXPOSE 8098
CMD ["./otp-service"]
