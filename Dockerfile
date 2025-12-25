# Production-ready Dockerfile for Streamify
# Build stage
FROM golang:1.25.5-alpine AS build
RUN apk add --no-cache curl
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o streamify cmd/api/main.go

# Final minimal image
FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/streamify /app/streamify

# Expose port (default 8080, can be overridden)
ENV PORT=8080
EXPOSE 8080

# Run the binary
CMD ["/app/streamify"]
