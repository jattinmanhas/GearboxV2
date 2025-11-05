#!/bin/bash

# Payment Service Migration Runner
# This script runs database migrations for the payment service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-}
DB_NAME=${DB_NAME:-gearbox_payment}
DOCKER_CONTAINER_NAME=${DOCKER_CONTAINER_NAME:postgres}
USE_DOCKER=${USE_DOCKER:true}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Set MIGRATIONS_DIR to the script's directory (where migration files are)
MIGRATIONS_DIR="$SCRIPT_DIR"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to find PostgreSQL Docker container
find_postgres_container() {
    if [ -n "$DOCKER_CONTAINER_NAME" ]; then
        echo "$DOCKER_CONTAINER_NAME"
        return
    fi
    
    # Try to find a postgres container
    local container=$(docker ps --format "{{.Names}}" | grep -i postgres | head -n 1)
    if [ -n "$container" ]; then
        echo "$container"
        return
    fi
    
    return 1
}

# Function to run migration
run_migration() {
    local migration_file=$1
    local direction=$2
    
    if [ ! -f "$migration_file" ]; then
        print_error "Migration file not found: $migration_file"
        exit 1
    fi
    
    print_status "Running migration: $(basename $migration_file) ($direction)"
    
    # Get absolute path of migration file for Docker volume mounting
    local migration_abs_path=$(cd "$(dirname "$migration_file")" && pwd)/$(basename "$migration_file")
    
    if [ "$USE_DOCKER" = "true" ] || [ -n "$DOCKER_CONTAINER_NAME" ]; then
        # Use Docker to run psql
        local container=$(find_postgres_container)
        if [ -z "$container" ]; then
            print_error "PostgreSQL Docker container not found. Set DOCKER_CONTAINER_NAME or ensure a postgres container is running."
            exit 1
        fi
        
        print_status "Using Docker container: $container"
        
        # Copy migration file to container and execute, or use docker exec with stdin
        docker exec -i "$container" psql -U "$DB_USER" -d "$DB_NAME" < "$migration_file"
    else
        # Use local psql
        if [ -n "$DB_PASSWORD" ]; then
            PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
        else
            psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
        fi
    fi
    
    if [ $? -eq 0 ]; then
        print_status "Migration completed successfully: $(basename $migration_file)"
    else
        print_error "Migration failed: $(basename $migration_file)"
        exit 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up      Run all pending migrations"
    echo "  down    Rollback the last migration"
    echo "  status  Show migration status"
    echo "  reset   Rollback all migrations and run them again"
    echo ""
    echo "Environment variables:"
    echo "  DB_HOST              Database host (default: localhost)"
    echo "  DB_PORT              Database port (default: 5432)"
    echo "  DB_USER              Database user (default: postgres)"
    echo "  DB_PASSWORD          Database password"
    echo "  DB_NAME              Database name (default: gearbox_payment)"
    echo "  USE_DOCKER           Use Docker to run psql (default: false)"
    echo "  DOCKER_CONTAINER_NAME PostgreSQL Docker container name (auto-detected if not set)"
    echo ""
    echo "Examples:"
    echo "  $0 up                                    # Run all migrations (local psql)"
    echo "  USE_DOCKER=true $0 up                   # Run migrations via Docker"
    echo "  DOCKER_CONTAINER_NAME=my-postgres $0 up # Use specific container"
    echo "  $0 down                                  # Rollback last migration"
    echo "  $0 status                                # Show migration status"
    echo "  DB_HOST=prod.example.com $0 up          # Run migrations on production"
}

# Parse command line arguments (needed for error messages)
COMMAND=${1:-up}

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    print_error "Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Check if psql is available locally or if Docker should be used
if ! command_exists psql; then
    if [ "$USE_DOCKER" = "true" ] || [ -n "$DOCKER_CONTAINER_NAME" ]; then
        print_warning "psql not found locally. Will use Docker to run migrations."
        USE_DOCKER=true
    elif command_exists docker; then
        # Try to auto-detect if we should use Docker
        if find_postgres_container > /dev/null 2>&1; then
            print_warning "psql not found locally. Auto-detected PostgreSQL Docker container. Using Docker mode."
            USE_DOCKER=true
        else
            print_error "psql command not found and no PostgreSQL Docker container detected."
            echo ""
            echo "Options:"
            echo "  1. Install PostgreSQL client: brew install postgresql@15"
            echo "  2. Use Docker: USE_DOCKER=true $0 $COMMAND"
            echo "  3. Set container name: DOCKER_CONTAINER_NAME=your-container $0 $COMMAND"
            exit 1
        fi
    else
        print_error "psql command not found and Docker is not available."
        echo ""
        echo "Please install PostgreSQL client tools:"
        echo "  brew install postgresql@15"
        exit 1
    fi
fi

case $COMMAND in
    up)
        print_status "Running migrations UP..."
        for migration in "$MIGRATIONS_DIR"/*.up.sql; do
            if [ -f "$migration" ]; then
                run_migration "$migration" "up"
            fi
        done
        print_status "All migrations completed successfully!"
        ;;
    down)
        print_warning "Rolling back migrations DOWN..."
        for migration in "$MIGRATIONS_DIR"/*.down.sql; do
            if [ -f "$migration" ]; then
                run_migration "$migration" "down"
            fi
        done
        print_status "All migrations rolled back successfully!"
        ;;
    status)
        print_status "Checking migration status..."
        echo "Database: $DB_NAME@$DB_HOST:$DB_PORT"
        echo "User: $DB_USER"
        echo "Migrations directory: $MIGRATIONS_DIR"
        echo ""
        echo "Available migration files:"
        ls -la "$MIGRATIONS_DIR"/*.sql 2>/dev/null || echo "No migration files found"
        ;;
    reset)
        print_warning "Resetting all migrations..."
        print_warning "This will DROP all tables and recreate them!"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Rolling back all migrations..."
            for migration in "$MIGRATIONS_DIR"/*.down.sql; do
                if [ -f "$migration" ]; then
                    run_migration "$migration" "down"
                fi
            done
            print_status "Running all migrations..."
            for migration in "$MIGRATIONS_DIR"/*.up.sql; do
                if [ -f "$migration" ]; then
                    run_migration "$migration" "up"
                fi
            done
            print_status "Database reset completed successfully!"
        else
            print_status "Reset cancelled."
        fi
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        echo ""
        show_usage
        exit 1
        ;;
esac
