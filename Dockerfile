# Production stage
FROM golang:alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

WORKDIR /src

# Create a structure that satisfies the 'replace' directive in go.mod
# github.com/halooid/backend/go-shared => ../go-shared
COPY backend/go-shared /src/backend/go-shared
COPY backend/auth-service /src/backend/auth-service

WORKDIR /src/backend/auth-service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service cmd/server/main.go

# Run stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/auth-service .

EXPOSE 50051

CMD ["./auth-service"]
