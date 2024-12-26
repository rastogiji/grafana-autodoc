# Stage 1: Build the Go app
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main .

# Stage 2: Create a lean production image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]