.PHONY: help build up down logs clean restart migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  make build      - Build all Docker images"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View logs from all services"
	@echo "  make clean      - Remove containers, volumes, and images"
	@echo "  make restart    - Restart all services"
	@echo "  make migrate    - Run database migrations"
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
	@echo "Note: Adjust these commands based on your migration setup"
	# docker-compose exec auth-service ./api migrate
	# docker-compose exec product-service ./api migrate
	# docker-compose exec payment-service ./api migrate
	# docker-compose exec blog-service npm run db:migrate

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
	docker-compose exec postgres-auth psql -U postgres -d gearbox_auth

db-product:
	docker-compose exec postgres-product psql -U postgres -d gearbox_product

db-payment:
	docker-compose exec postgres-payment psql -U postgres -d gearbox_payment

db-blog:
	docker-compose exec postgres-blog psql -U postgres -d gearbox_blog

