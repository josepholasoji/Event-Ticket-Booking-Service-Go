# ============================
# Build stage
# ============================
FROM golang:1.25.6-alpine AS builder

# Install system dependencies required by Buffalo
RUN apk add --no-cache \
	git \
	bash \
	build-base \
	postgresql-client

# Install Buffalo CLI
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest

WORKDIR /app

# Copy Go module files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy configuration files required by Pop
COPY config/ ./config/
COPY database.yml ./
# Copy the rest of the application
COPY . .

# Build the Buffalo app for production
# -static: statically link libc to avoid runtime dependencies
# -o bin/app: output to bin directory
RUN buffalo build -v --environment production --static -o bin/app

# ============================
# Runtime stage
# ============================
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache \
	ca-certificates \
	postgresql-client \
	bash

WORKDIR /app

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

# Copy the compiled binary from builder
COPY --from=builder --chown=app:app /app/bin/app /app/app

# Copy configuration files and migrations
COPY --from=builder --chown=app:app /app/config /app/config
COPY --from=builder --chown=app:app /app/database.yml /app/database.yml
COPY --from=builder --chown=app:app /app/migrations /app/migrations

# Set runtime user
USER app

# Expose Buffalo default port
EXPOSE 3000

# Optional runtime environment defaults
ENV GO_ENV=production

# Lightweight healthcheck (returns 0 if binary exists)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
	CMD ["/bin/sh", "-c", "[ -x /app/app ] && echo ok || exit 1"]

# Run the application
ENTRYPOINT ["/app/app"]

