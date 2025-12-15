FRONTEND_DIR = ./web
BACKEND_DIR = .

.PHONY: all build-frontend start-backend dev-db dev-db-down dev-db-logs

all: build-frontend start-backend

# Development database services
dev-db:
	@echo "Starting development database services..."
	@docker-compose -f docker-compose.dev.yml up -d

dev-db-down:
	@echo "Stopping development database services..."
	@docker-compose -f docker-compose.dev.yml down

dev-db-logs:
	@docker-compose -f docker-compose.dev.yml logs -f

build-frontend:
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && bun install && DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat VERSION) bun run build

start-backend:
	@echo "Starting backend dev server..."
	@cd $(BACKEND_DIR) && go run main.go &
