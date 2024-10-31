# Build stage
FROM golang:1.22-alpine AS builder

# Build args
ARG VERSION="dev"

WORKDIR /app
COPY . .
RUN go mod download

# Build with version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X main.version=${VERSION}" \
    -o nanodns ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/nanodns .
EXPOSE 53/udp
CMD ["./nanodns"]