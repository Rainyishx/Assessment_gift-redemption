#Stage 1: Build 
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go module files first for better Docker layer caching
COPY go.mod ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o gift-redemption ./cmd/main.go

# Stage 2: Run 
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/gift-redemption .

# Copy the sample data directory
COPY data/ ./data/

# Expose the default port
EXPOSE 8080

# Environment variable defaults 
ENV PORT=8080
ENV STAFF_MAPPING_FILE=data/staffmapping.csv
ENV REDEMPTION_FILE=data/redemptions.csv

CMD ["./gift-redemption"]