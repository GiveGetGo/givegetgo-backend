#!/bin/bash

# Define repository details
USERNAME="givegetgo"
REPOSITORY="givegetgo-backend"

# Authenticate to GitHub Packages
echo $PAT | docker login ghcr.io -u $USERNAME --password-stdin

# Build all services
docker-compose build

# Check build status
if [ $? -ne 0 ]; then
  echo "Build failed. Exiting."
  exit 1
fi

# Declare an array of the full names of your services
services=("ghcr.io/$USERNAME/$REPOSITORY/givegetgo-user-backend:latest" 
          "ghcr.io/$USERNAME/$REPOSITORY/givegetgo-verification-backend:latest" 
          "ghcr.io/$USERNAME/$REPOSITORY/givegetgo-post-backend:latest" 
          "ghcr.io/$USERNAME/$REPOSITORY/givegetgo-bid-backend:latest" 
          "ghcr.io/$USERNAME/$REPOSITORY/givegetgo-match-backend:latest" 
          "ghcr.io/$USERNAME/$REPOSITORY/givegetgo-notification-backend:latest")

# Push each service
for service in "${services[@]}"
do  
  echo "Pushing $service"
  docker push $service

  # Check if push was successful
  if [ $? -ne 0 ]; then
    echo "Failed to push $service. Exiting."
    exit 1
  fi
done

echo "All images have been pushed successfully."
