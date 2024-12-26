# Stage 1: Build the binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main .

# Stage 2: Cpoy Binary to alpine image to create a lean image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
RUN chmod +x /root/main

ENTRYPOINT ["/root/main"]