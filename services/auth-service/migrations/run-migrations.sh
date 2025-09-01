#!/bin/bash

# =====================================================
# Migration Runner Script for Auth Service
# =====================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
MIGRATION_PATH="./migrations"
DATABASE_URL=""
COMMAND=""

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] COMMAND"
    echo ""
    echo "Commands:"
    echo "  up        Apply all pending migrations"
    echo "  down      Rollback last migration"
    echo "  down-all  Rollback all migrations"
    echo "  version   Show current migration version"
    echo "  force     Force migration to specific version"
    echo "  create    Create new migration files"
    echo ""
    echo "Options:"
    echo "  -d, --database DATABASE_URL    Database connection string"
    echo "  -p, --path PATH               Migration files path (default: ./migrations)"
    echo "  -h, --help                    Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -d 'postgres://user:pass@localhost:5432/dbname?sslmode=disable' up"
    echo "  $0 -d 'postgres://user:pass@localhost:5432/dbname?sslmode=disable' down 1"
    echo "  $0 -d 'postgres://user:pass@localhost:5432/dbname?sslmode=disable' version"
}

# Function to check if go-migrate is installed
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        print_error "go-migrate is not installed. Please install it first:"
        echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

# Function to validate database URL
validate_database_url() {
    if [[ -z "$DATABASE_URL" ]]; then
        print_error "Database URL is required. Use -d or --database option."
        exit 1
    fi
}

# Function to check if migrations directory exists
check_migrations_dir() {
    if [[ ! -d "$MIGRATION_PATH" ]]; then
        print_error "Migrations directory not found: $MIGRATION_PATH"
        exit 1
    fi
}

# Function to run migration command
run_migration() {
    local cmd="$1"
    local args="$2"
    
    print_status "Running migration: migrate -path $MIGRATION_PATH -database \"$DATABASE_URL\" $cmd $args"
    
    if migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" "$cmd" $args; then
        print_success "Migration completed successfully"
    else
        print_error "Migration failed"
        exit 1
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--database)
            DATABASE_URL="$2"
            shift 2
            ;;
        -p|--path)
            MIGRATION_PATH="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        up|down|down-all|version|force|create)
            COMMAND="$1"
            shift
            ;;
        *)
            # Additional arguments for commands like 'down 1' or 'force 2'
            ARGS="$ARGS $1"
            shift
            ;;
    esac
done

# Validate inputs
if [[ -z "$COMMAND" ]]; then
    print_error "Command is required"
    show_usage
    exit 1
fi

# Check prerequisites
check_migrate
validate_database_url
check_migrations_dir

# Execute command
case $COMMAND in
    up)
        print_status "Applying all pending migrations..."
        run_migration "up"
        ;;
    down)
        if [[ -n "$ARGS" ]]; then
            print_status "Rolling back $ARGS migration(s)..."
            run_migration "down" "$ARGS"
        else
            print_status "Rolling back last migration..."
            run_migration "down" "1"
        fi
        ;;
    down-all)
        print_warning "Rolling back ALL migrations. This will remove all data!"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Rolling back all migrations..."
            run_migration "down"
        else
            print_status "Operation cancelled"
            exit 0
        fi
        ;;
    version)
        print_status "Checking current migration version..."
        run_migration "version"
        ;;
    force)
        if [[ -z "$ARGS" ]]; then
            print_error "Version number required for force command"
            exit 1
        fi
        print_warning "Forcing migration to version $ARGS. This may cause data loss!"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Forcing migration to version $ARGS..."
            run_migration "force" "$ARGS"
        else
            print_status "Operation cancelled"
            exit 0
        fi
        ;;
    create)
        if [[ -z "$ARGS" ]]; then
            print_error "Migration name required for create command"
            exit 1
        fi
        print_status "Creating new migration: $ARGS"
        migrate create -ext sql -dir "$MIGRATION_PATH" "$ARGS"
        print_success "Migration files created successfully"
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        show_usage
        exit 1
        ;;
esac
