FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main .

RUN ls -l /app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

RUN ls -l /root
RUN chmod +x /root/main

ENTRYPOINT ["./main"]