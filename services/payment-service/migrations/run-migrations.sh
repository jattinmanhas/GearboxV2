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
DB_NAME=${DB_NAME:-gearbox_payment}
MIGRATIONS_DIR="./migrations"

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

# Function to run migration
run_migration() {
    local migration_file=$1
    local direction=$2
    
    if [ ! -f "$migration_file" ]; then
        print_error "Migration file not found: $migration_file"
        exit 1
    fi
    
    print_status "Running migration: $(basename $migration_file) ($direction)"
    
    if [ "$direction" = "up" ]; then
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
    else
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
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
    echo "  DB_HOST     Database host (default: localhost)"
    echo "  DB_PORT     Database port (default: 5432)"
    echo "  DB_USER     Database user (default: postgres)"
    echo "  DB_NAME     Database name (default: gearbox_payment)"
    echo ""
    echo "Examples:"
    echo "  $0 up                    # Run all migrations"
    echo "  $0 down                  # Rollback last migration"
    echo "  $0 status                # Show migration status"
    echo "  DB_HOST=prod.example.com $0 up  # Run migrations on production"
}

# Check if psql is available
if ! command_exists psql; then
    print_error "psql command not found. Please install PostgreSQL client tools."
    exit 1
fi

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    print_error "Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Parse command line arguments
COMMAND=${1:-up}

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
