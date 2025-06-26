FROM redis:7-alpine

# Create a directory for Redis data
RUN mkdir -p /data

# Expose Redis port
EXPOSE 6379

# Use default Redis configuration
CMD ["redis-server", "--appendonly", "yes"] 