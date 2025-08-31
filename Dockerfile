FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN adduser -D -s /bin/sh autodoc

WORKDIR /home/autodoc

COPY --from=builder /app/cmd/bin/autodoc ./grafana-autodoc

RUN chmod +x ./grafana-autodoc && \
    chown autodoc:autodoc ./grafana-autodoc

USER autodoc

ENTRYPOINT ["./grafana-autodoc"]