# Lightweight base image for building
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend cmd/server/*.go

# Minimal base image for deployment
FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /app/backend .

COPY .env .

CMD ["./backend"]