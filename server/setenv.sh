#!/usr/bin/env bash
# Configure Postgres connection environment variables for the Go API

# Adjust these if your container uses different values
export PGHOST=postgres_vector          # or postgres_vector if inside Docker network
export PGPORT=5432
export PGUSER=user
export PGPASSWORD=password
export PGDATABASE=rag_db
export PGSSLMODE=disable

echo "✅ Environment variables set for Postgres connection"
echo "   Host: $PGHOST"
echo "   Port: $PGPORT"
echo "   DB:   $PGDATABASE"

