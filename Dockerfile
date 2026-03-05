# --- development stage ---
FROM golang:1.26-alpine AS development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod .
RUN go mod download
COPY . .
CMD ["air"]

# --- builder stage ---
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# --- production stage ---
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
