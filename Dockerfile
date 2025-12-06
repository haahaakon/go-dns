# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-dns .

# Final stage
FROM alpine:latest


WORKDIR /app/

# Copy the binary from builder
COPY --from=builder /app/go-dns .

EXPOSE 5353/udp

CMD ["./go-dns"]
