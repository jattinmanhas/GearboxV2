-- Database initialization script for Gearbox Services
-- This script creates the necessary databases for all microservices
-- It is meant to be run manually or during initial Docker container startup

-- Create gearbox_auth database
SELECT 'CREATE DATABASE gearbox_auth'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gearbox_auth')\gexec

-- Create gearbox_product database
SELECT 'CREATE DATABASE gearbox_product'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gearbox_product')\gexec

-- Create gearbox_payment database
SELECT 'CREATE DATABASE gearbox_payment'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gearbox_payment')\gexec

-- Create gearbox_blog database
SELECT 'CREATE DATABASE gearbox_blog'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gearbox_blog')\gexec
