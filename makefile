# Variables
ENV_FILE = .env
ENV_SAMPLE_FILE = .env.example
DOCKER_COMPOSE = docker compose
DOCKER_COMPOSE_CMD = $(DOCKER_COMPOSE) up --build

all: copy-env up 

copy-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		cp $(ENV_SAMPLE_FILE) $(ENV_FILE); \
		echo "$(ENV_FILE) created from $(ENV_SAMPLE_FILE)"; \
	else \
		echo "$(ENV_FILE) already exists"; \
	fi

# Start Docker containers
up:
	$(DOCKER_COMPOSE_CMD)

# Stop Docker containers
down:
	$(DOCKER_COMPOSE) down

# Rebuild Docker images and start containers
restart: down up

# Remove all stopped containers and dangling images
clean:
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

test:
	go test -v ./...

mockstore:
	mockgen -package mocks -destination mocks/store.go github.com/kenmobility/git-api-service/test Store

mockgit:
	mockgen -package mocks -destination mocks/mock_git_manager_client.go github.com/kenmobility/git-api-service/infra/git GitManagerClient

.PHONY: all copy-env up down restart clean test mockstore mockgit
