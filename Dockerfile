FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Build the main server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o aviation-service ./cmd/server/main.go
# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go
# Build the cron job tool
RUN CGO_ENABLED=0 GOOS=linux go build -o schedule ./cmd/schedule/main.go

# Create a minimal production image
FROM alpine:3.22

WORKDIR /app

# Copy binaries from the builder stage
COPY --from=builder /app/aviation-service /app/aviation-service
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/schedule /app/schedule

# Copy migrations
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/seedings /app/seedings

COPY .env .env

# Set the entry point to the main server by default
ENTRYPOINT ["/app/aviation-service"]
