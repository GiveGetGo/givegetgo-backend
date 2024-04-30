#!/bin/bash

# Path to the directory containing service directories
base_dir="./servers"

# Iterate over each service directory
for service_dir in "$base_dir"/*; do
    if [ -d "$service_dir" ]; then  # Check if it's a directory
        service_name=$(basename "$service_dir")
        env_file="$service_dir/.env.$service_name"

        if [ -f "$env_file" ]; then  # Check if the .env file exists
            # Define the output file name
            output_file="${service_name}-encoded.txt"

            # Encode the file to base64 and save it to the output file
            echo "Encoding $env_file into $output_file"
            base64 -i "$env_file" > "$output_file"
        else
            echo "No .env file found for $service_name"
        fi
    fi
done

echo "Encoding completed. Encoded files are in the current directory."
