#!/bin/sh

# Use the REDIS_PASSWORD environment variable to set the password for Redis
if [ -n "$REDIS_PASSWORD" ]; then
    set -- redis-server --requirepass "$REDIS_PASSWORD"
else
    set -- redis-server
fi

# Execute the Redis server with any provided arguments
exec "$@"
