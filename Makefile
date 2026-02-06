ROOT := .
TMP_DIR := tmp
APP_NAME := news_article_app
BIN := $(TMP_DIR)/main
CSS_INPUT := assets/css/main.css
CSS_OUTPUT := static/css/styles.css

EXCLUDE_DIRS := tmp node_modules vendor
EXCLUDE_REGEX := .*_templ.go _test.go
INCLUDE_DIRS := templates assets handlers models repositories services
INCLUDE_FILES := main.go
INCLUDE_EXTS := go templ html css

# Database Configuration (PLACEHOLDERS ONLY - DO NOT COMMIT REAL VALUES)
DB_NAME := DB_NAME
DB_USER := DB_USER
DB_PASS := DB_PASSWORD
PG_ADMIN := PG_ADMIN_USER
PG_ADMIN_PASS := PG_ADMIN_PASSWORD

PSQL := env PGPASSWORD=$(PG_ADMIN_PASS) psql -U $(PG_ADMIN) -h localhost

.PHONY: all build clean watch run install stop

all: build

install:
	@echo "Installing Go dependencies..."
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download

	@echo "Installing TailwindCSS..."
	@npm install -D tailwindcss

	@echo "Installing PostgreSQL..."
	@if [ -f "/etc/debian_version" ]; then \
		sudo apt update && sudo apt install -y postgresql postgresql-contrib; \
	elif [ -f "/etc/alpine-release" ]; then \
		sudo apk update && sudo apk add postgresql postgresql-contrib; \
	else \
		echo "Unsupported OS. Please install PostgreSQL manually."; \
	fi

build: clean
	@echo "Building project..."
	templ generate templates
	go build -o $(APP_NAME) ./main.go
	npx @tailwindcss/cli -i $(CSS_INPUT) -o $(CSS_OUTPUT)

clean:
	@echo "Cleaning up..."
	rm -rf $(TMP_DIR)
	mkdir -p $(TMP_DIR)
	rm -f $(APP_NAME) $(CSS_OUTPUT)

watch:
	@echo "Starting file watcher..."
	@find $(INCLUDE_DIRS) $(INCLUDE_FILES) -type f \( $(foreach ext,$(INCLUDE_EXTS),-name '*.$(ext)') \) | entr -r make build

run: build
	@echo "Running application..."
	@./$(APP_NAME)

stop:
	@echo "Stopping application..."
	@pkill -f $(APP_NAME) || true

# -------- CREATE SECTION --------
db-create:
	@echo "Creating PostgreSQL user: $(DB_USER)..."
	@$(PSQL) -c "CREATE USER $(DB_USER) WITH PASSWORD '$(DB_PASS)';"
	@echo "Creating database: $(DB_NAME)..."
	@$(PSQL) -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"
	@echo "Granting privileges to $(DB_USER)..."
	@$(PSQL) -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_NAME) TO $(DB_USER);"
	@PGPASSWORD="$(DB_PASS)" psql -U $(DB_USER) -d $(DB_NAME) -h localhost -f postgres_db/schema.sql

db-insert-admin:
	@echo "Inserting admin user into $(DB_NAME)..."
	@PGPASSWORD=$(DB_PASS) psql -U $(DB_USER) -d $(DB_NAME) -h localhost -c \
	"INSERT INTO users (username, email, password_hash, role, created_at) \
	VALUES ('admin', 'admin@placeholder.invalid', 'REPLACE_WITH_HASHED_PASSWORD', 'admin', NOW()) \
	ON CONFLICT (email) DO NOTHING;"

db-insert-mock-article:
	@echo "Inserting mock article into $(DB_NAME)..."
	@PGPASSWORD=$(DB_PASS) psql -U $(DB_USER) -d $(DB_NAME) -h localhost -c \
	"INSERT INTO articles (title, content, author_id, created_at, updated_at) \
	VALUES ('Mock Article Title', 'This is a test article for development purposes.', 1, NOW(), NOW()) \
	ON CONFLICT (id) DO NOTHING;"


# -------- DELETE SECTION --------
db-delete:
	@echo "Dropping database: $(DB_NAME)..."
	@$(PSQL) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Dropping user: $(DB_USER)..."
	@$(PSQL) -c "DROP USER IF EXISTS $(DB_USER);"
