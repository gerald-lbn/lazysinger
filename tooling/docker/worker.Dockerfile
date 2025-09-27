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
RUN go build ./cmd/worker

FROM base
WORKDIR /app
COPY --from=builder /app/worker /app/worker

ENTRYPOINT [ "/app/worker" ]
