FROM golang:1.24.9-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app
VOLUME [ "/data", "/music" ]
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config/config.yml /app/config/config.yml
RUN chmod +x /app/main
ENV REFRAIN_APP_ENVIRONMENT=production
ENTRYPOINT [ "/app/main" ]
