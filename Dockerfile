FROM golang:1.24-alpine AS base

FROM base AS builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o lazysinger ./cmd

FROM base
WORKDIR /app
COPY --from=builder /app/lazysinger /app/lazysinger
RUN chmod +x /app/lazysinger
ENTRYPOINT ["/app/lazysinger"]
