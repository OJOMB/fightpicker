# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

# Enable static build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install CA certs for HTTPS clients
RUN apk add --no-cache ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o server \
    ./cmd/server

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy certs & binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /app/server

# Non-root user
USER nonroot:nonroot

# Expose app port (adjust if needed)
EXPOSE 8080

ENTRYPOINT ["/app/server"]