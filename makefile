FRONTEND_DIR = ./web
BACKEND_DIR = .

.PHONY: all build-frontend start-backend dev-db dev-db-down dev-db-logs dev dev-stop

all: build-frontend start-backend

# Full development environment - runs all services with logs in one console
# Usage: make dev
# Stop with: Ctrl+C (or make dev-stop in another terminal)
dev:
	@echo "Starting full development environment..."
	@echo "  - Docker services (PostgreSQL, Redis)"
	@echo "  - Backend (Go on port 3000)"
	@echo "  - Frontend (Vite on port 5173)"
	@echo ""
	@echo "Press Ctrl+C to stop all services"
	@echo "================================================"
	@docker-compose -f docker-compose.dev.yml up -d
	@sleep 2
	@(trap 'kill 0' SIGINT; \
		(cd $(BACKEND_DIR) && go run main.go 2>&1 | sed 's/^/[backend] /') & \
		(cd $(FRONTEND_DIR) && bun run dev 2>&1 | sed 's/^/[frontend] /') & \
		(docker-compose -f docker-compose.dev.yml logs -f 2>&1 | sed 's/^/[docker] /') & \
		wait)

# Stop all development services
dev-stop:
	@echo "Stopping all development services..."
	@docker-compose -f docker-compose.dev.yml down
	@-pkill -f "go run main.go" 2>/dev/null || true
	@-pkill -f "vite" 2>/dev/null || true
	@echo "All services stopped."

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
