# --- Development/Debug Stage ---
FROM golang:1.22-alpine AS dev-stage
WORKDIR /app
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# --- Build Stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o out/main ./cmd

# --- Final Stage ---
FROM alpine:3.20
COPY --from=builder /app/out/main /app/main
EXPOSE 8080
CMD ["/app/main"]