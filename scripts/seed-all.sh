#!/bin/sh

set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

seed_sql() {
  service_name="postgres"
  database_name="$2"
  seed_file="$3"

  printf 'Seeding %s...\n' "$database_name"
  docker-compose exec -T "$service_name" psql -U postgres -d "$database_name" < "$ROOT_DIR/$seed_file"
}

seed_sql postgres gearbox_auth services/auth-service/seeds/001_demo_seed.sql
seed_sql postgres gearbox_product services/product-service/seeds/001_demo_seed.sql
seed_sql postgres gearbox_payment services/payment-service/seeds/001_demo_seed.sql
seed_sql postgres gearbox_blog services/blog-service/seeds/001_demo_seed.sql

printf '\nSeeded all local Gearbox databases.\n'
printf 'Demo login: ava.admin@gearbox.local / Password123!\n'
