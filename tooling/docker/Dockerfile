FROM golang:1.24-alpine AS base

FROM base AS builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o lazysinger ./cmd

FROM base
WORKDIR /app
RUN addgroup --system appuser
RUN adduser --system appuser
COPY --from=builder --chown=appuser:appuser /app/lazysinger /app/lazysinger
RUN chmod +x /app/lazysinger
USER appuser

ENTRYPOINT ["/app/lazysinger"]
