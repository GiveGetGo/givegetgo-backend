# Start from a Debian-based image with the Go SDK
FROM golang:1.22.0 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./
# Download any necessary dependencies
RUN go mod download

# Copy the rest of the application's code
COPY . .

# Compile the application to /main.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a Docker multi-stage build to create a lean image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the .env file specific to the verification service
COPY .env.user /root/.env.user

# Copy the configuration directory
COPY --from=builder /app/config /root/config

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
