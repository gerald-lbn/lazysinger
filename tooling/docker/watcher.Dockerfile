FROM golang:1.24-alpine AS base

FROM base AS builder
WORKDIR /app
# ensure a portable and static-ish binart
ENV CGO_ENABLE=0
ENV GOOS=linux
ENV COARCH=amd64
# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
RUN go build ./cmd/watcher

FROM base
WORKDIR /app
# Create non-root user
RUN addgroup --system appuser
RUN adduser --system appuser
# Copy binary and change ownership
COPY --from=builder --chown=appuser:appuser /app/watcher /app/watcher
# Run as non-root
USER appuser

ENTRYPOINT [ "/app/watcher" ]
