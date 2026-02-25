# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download deps first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ssh-portal .

# ── Run stage ─────────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache openssh-keygen ca-certificates

WORKDIR /app

COPY --from=builder /app/ssh-portal .

# Create data dir for host key persistence
RUN mkdir -p /app/data/.ssh

# Generate host key at startup if not mounted
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 2222

ENTRYPOINT ["/app/entrypoint.sh"]
