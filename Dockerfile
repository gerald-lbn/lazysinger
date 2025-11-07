FROM golang:1.24.9-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main ./main.go

FROM alpine:latest
WORKDIR /app
VOLUME [ "/music" ]
COPY --from=builder /app/main /app
RUN chmod +x /app/main
ENTRYPOINT [ "/app/main" ]
