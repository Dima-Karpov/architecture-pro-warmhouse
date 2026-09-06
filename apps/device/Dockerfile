# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o device ./cmd/api

# Runtime stage
FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata wget \
    && addgroup -g 1000 app \
    && adduser -D -u 1000 -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /app/device .

USER app

EXPOSE 8082

CMD ["./device"]
