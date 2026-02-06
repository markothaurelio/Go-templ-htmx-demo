# Stage 1: Build frontend assets using Node.js
FROM node:latest AS frontend

WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy Tailwind config and source files
COPY tailwind.config.js ./
COPY assets/css/main.css ./assets/css/main.css
COPY templates/ ./templates/

# Ensure Tailwind scans all necessary files
RUN npx tailwindcss -i assets/css/main.css -o static/css/styles.css --config tailwind.config.js

# Stage 2: Build the Go application
FROM golang:1.23 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install Templ (Fixes "templ: not found" error)
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire source code
COPY . .

# Copy the generated static assets from frontend build
COPY --from=frontend /app/static ./static

# Run Templ to generate Go templates and build the binary
RUN templ generate templates && go build -o news_article_app ./main.go

# Stage 3: Create a lightweight production image with PostgreSQL and DB seeding
FROM alpine:latest

# Install certificates and PostgreSQL (server + client)
RUN apk --no-cache add ca-certificates postgresql postgresql-contrib su-exec

# Create /run/postgresql directory and set permissions
RUN mkdir -p /run/postgresql && chown postgres:postgres /run/postgresql

WORKDIR /root/

# Environment variables for the Go app and database configuration (PLACEHOLDERS ONLY)
ENV GO_EMAIL="" \
    GO_EMAIL_PASS="" \
    JWT_KEY="REPLACE_WITH_SECURE_VALUE" \
    DB_NAME="DB_NAME" \
    DB_USER="DB_USER" \
    DB_PASS="DB_PASSWORD" \
    PG_ADMIN="PG_ADMIN_USER" \
    PG_ADMIN_PASS="PG_ADMIN_PASSWORD"

# Copy the compiled binary and static assets
COPY --from=builder /app/news_article_app .
COPY --from=builder /app/static /root/static

# Copy the geojson folder into the same dir as the binary
COPY geojson_data /root/geojson_data

# Copy the database schema (assumes your schema file is at postgres_db/schema.sql)
COPY postgres_db /root/postgres_db

# Copy the entrypoint script that initializes PostgreSQL and seeds the DB
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the application port
EXPOSE 3000

# Run the entrypoint script
CMD ["/entrypoint.sh"]
