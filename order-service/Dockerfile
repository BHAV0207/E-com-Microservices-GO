# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup \
  && adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

COPY --from=builder /app/main .

RUN chown -R appuser:appgroup /root/
USER appuser

# Note: order-service doesn't expose an HTTP port yet
# EXPOSE 7000

CMD ["./main"]


