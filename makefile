# List of services
SERVICES = user bid match post verification notification

# Base directory for servers
SERVERS_DIR = ./servers

# Default target
all: env_files up

# Generate .env files if they don't exist
env_files:
	@$(foreach service, $(SERVICES), \
		if [ ! -f $(SERVERS_DIR)/$(service)/.env.$(service) ]; then \
			echo "Generating $(SERVERS_DIR)/$(service)/.env.$(service) from $(SERVERS_DIR)/$(service)/.env.$(service).example..."; \
			cp $(SERVERS_DIR)/$(service)/.env.$(service).example $(SERVERS_DIR)/$(service)/.env.$(service); \
		else \
			echo "$(SERVERS_DIR)/$(service)/.env.$(service) already exists."; \
		fi;)

.env:
	@if [ ! -f .env ]; then \
		echo "Generating .env from .env.example..."; \
		cp .env.example .env; \
	else \
		echo ".env already exists."; \
	fi

redis/.env.redis:
	@if [ ! -f redis/.env.redis ]; then \
		echo "Generating redis/.env.redis from redis/.env.redis.example..."; \
		cp redis/.env.redis.example redis/.env.redis; \
	else \
		echo "redis/.env.redis already exists."; \
	fi

# Bring up the project using Docker Compose
up:
	@echo "Starting Docker Compose..."
	docker-compose up -d --build

# Clean up the environment files (optional)
clean:
	@echo "Cleaning up environment files..."
	rm -f .env redis/.env.redis
	@$(foreach service, $(SERVICES), rm -f $(SERVERS_DIR)/$(service)/.env.$(service);)
