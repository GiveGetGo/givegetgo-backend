# Use the official Redis image as the base
FROM redis:latest

# Copy the custom entrypoint script to the container
COPY entrypoint.sh /usr/local/bin/

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint script to be executed
ENTRYPOINT ["entrypoint.sh"]
