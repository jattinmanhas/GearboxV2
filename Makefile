.PHONY: help build up down logs clean restart migrate seed db-up db-down local-all local-stop local-auth local-product local-payment local-blog local-frontend

# Default target
help:
	@echo "Available commands:"
	@echo "  make build      - Build all Docker images"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make db-up      - Start only the PostgreSQL containers"
	@echo "  make db-down    - Stop only the PostgreSQL containers"
	@echo "  make logs       - View logs from all services"
	@echo "  make clean      - Remove containers, volumes, and images"
	@echo "  make restart    - Restart all services"
	@echo "  make migrate    - Run database migrations"
	@echo "  make seed       - Seed all local databases with demo data"
	@echo "  make local-all     - Run all services locally against Docker Postgres"
	@echo "  make local-auth     - Run auth-service locally against Docker Postgres"
	@echo "  make local-product  - Run product-service locally against Docker Postgres"
	@echo "  make local-payment  - Run payment-service locally against Docker Postgres"
	@echo "  make local-blog     - Run blog-service locally against Docker Postgres"
	@echo "  make local-frontend - Run frontend locally against local services"
	@echo "  make prod-up    - Start production services"
	@echo "  make prod-down  - Stop production services"

# Build all images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d
	@echo "Services starting... Check logs with: make logs"

# Stop all services
down:
	docker-compose down

# Start only postgres container for local app testing
db-up:
	docker-compose up -d postgres
	@echo "PostgreSQL container is starting on port 5432"

# Stop only postgres container
db-down:
	docker-compose stop postgres

# View logs
logs:
	docker-compose logs -f

# Clean everything
clean:
	docker-compose down -v
	docker system prune -f

# Restart services
restart:
	docker-compose restart

# Run migrations (adjust based on your migration setup)
migrate:
	@echo "Running migrations..."
	@echo "Note: Ensure your local services are configured to use the single Postgres on port 5432"

# Seed local databases
seed:
	sh scripts/seed-all.sh

local-stop:
	@echo "Stopping local services..."
	-lsof -ti :8081,8082,8083,3000,3001 | xargs kill -9 2>/dev/null || true
	@echo "Local services stopped."

# Run all services locally while keeping Postgres in Docker
local-all: local-stop
	@echo "Starting all services locally... Press Ctrl+C to stop."
	@trap 'echo "Stopping services..."; $(MAKE) local-stop; exit 0' INT TERM QUIT; \
	(cd services/auth-service && go run cmd/api/main.go) & \
	(cd services/product-service && go run cmd/api/main.go) & \
	(cd services/payment-service && go run cmd/api/main.go) & \
	(cd services/blog-service && npm run dev) & \
	(cd frontend && npm run dev) & \
	wait

# Run services locally while keeping Postgres in Docker
local-auth:
	cd services/auth-service && go run cmd/api/main.go

local-product:
	cd services/product-service && go run cmd/api/main.go

local-payment:
	cd services/payment-service && go run cmd/api/main.go

local-blog:
	cd services/blog-service && npm run dev

local-frontend:
	cd frontend && npm run dev

# Production commands
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

prod-down:
	docker-compose -f docker-compose.prod.yml down

# Build individual services
build-auth:
	docker build -f services/auth-service/Dockerfile -t gearbox-auth-service .

build-product:
	docker build -f services/product-service/Dockerfile -t gearbox-product-service .

build-payment:
	docker build -f services/payment-service/Dockerfile -t gearbox-payment-service .

build-blog:
	docker build -f services/blog-service/Dockerfile -t gearbox-blog-service .

build-frontend:
	docker build -f frontend/Dockerfile -t gearbox-frontend .

# View service status
status:
	docker-compose ps

# Access database
db-auth:
	docker-compose exec postgres psql -U postgres -d gearbox_auth

db-product:
	docker-compose exec postgres psql -U postgres -d gearbox_product

db-payment:
	docker-compose exec postgres psql -U postgres -d gearbox_payment

db-blog:
	docker-compose exec postgres psql -U postgres -d gearbox_blog
