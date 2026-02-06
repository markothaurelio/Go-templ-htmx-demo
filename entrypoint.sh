#!/bin/sh
set -e

# --- Initialize PostgreSQL ---
DATA_DIR="/var/lib/postgresql/data"
if [ ! -d "$DATA_DIR" ]; then
  mkdir -p "$DATA_DIR"
  chown -R postgres:postgres "$DATA_DIR"
  su postgres -c "initdb -D $DATA_DIR"
fi

# Start PostgreSQL in the background
su postgres -c "pg_ctl -D $DATA_DIR -o \"-c listen_addresses='localhost'\" -w start"

# Allow PostgreSQL a moment to fully start up
sleep 2

# --- Database Setup ---
# Create the database user
psql -U "$PG_ADMIN" -h localhost -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';" || true

# Create the database owned by the new user
psql -U "$PG_ADMIN" -h localhost -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" || true

# Grant privileges (if needed)
psql -U "$PG_ADMIN" -h localhost -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" || true

# Apply the schema from file (if not already applied)
PGPASSWORD="$DB_PASS" psql -U "$DB_USER" -d "$DB_NAME" -h localhost -f /root/postgres_db/schema.sql || true

# Seed the database: insert an admin user if not exists (PLACEHOLDER VALUES)
PGPASSWORD="$DB_PASS" psql -U "$DB_USER" -d "$DB_NAME" -h localhost -c "\
  INSERT INTO users (username, email, password_hash, role, created_at) \
  VALUES ('admin', 'admin@placeholder.invalid', 'REPLACE_WITH_HASHED_PASSWORD', 'admin', NOW()) \
  ON CONFLICT (email) DO NOTHING;"

# (Optional) List the static assets to verify they were copied correctly
ls -l /root/static/css

# --- Start the Go Application ---
exec ./news_article_app
