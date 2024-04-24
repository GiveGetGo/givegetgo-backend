# Use the official nginx image as a parent image
FROM nginx:latest

# Remove the default nginx configuration file
RUN rm /etc/nginx/conf.d/default.conf

# Copy the configuration file from the current directory to the container
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80
