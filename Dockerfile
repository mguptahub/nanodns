# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o nanodns ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/nanodns .
EXPOSE 53/udp
CMD ["./nanodns"]