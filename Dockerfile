# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./

COPY . .
# Build with security flags
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main.out ./cmd/


FROM alpine:3.20 AS runtime

# Install only the needed packages
RUN apk add --no-cache nginx

COPY nginx.conf /etc/nginx/http.d/default.conf

# Create non-root user
RUN adduser -D -g '' appuser && \
    # Remove setuid/setgid binaries
    find / -perm /4000 -exec chmod u-s {} \; || true && \
    find / -perm /2000 -exec chmod g-s {} \; || true && \
    # Create necessary directories with proper permissions
    mkdir -p /var/lib/nginx /var/log/nginx /run/nginx && \
    chown -R appuser:appuser /var/lib/nginx /var/log/nginx /run/nginx && \
    chmod -R 755 /var/lib/nginx /var/log/nginx /run/nginx
    
# Copy the compiled application from the builder stage
COPY --from=builder /app/main.out /usr/local/bin/main.out
RUN chmod +x /usr/local/bin/main.out
    
# EXPOSE 8080

# Switch to non-root user
USER appuser

# Start app and nginx
CMD ["/bin/sh", "-c", "/usr/local/bin/main.out --port=9999 & nginx -g 'daemon off;'"]