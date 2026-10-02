# Multi-stage Docker build for SpecAuditor Gov-Edition
# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go module files
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Compile self-contained static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o spec-auditor cmd/api/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/spec-auditor /app/spec-auditor

# Expose HTTP port
EXPOSE 8080

# Environment defaults
ENV PORT=8080 \
    OLLAMA_URL=http://ollama:11434 \
    OLLAMA_MODEL=qwen2.5:7b

ENTRYPOINT ["/app/spec-auditor"]
