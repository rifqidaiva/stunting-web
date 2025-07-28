# Makefile for building and running a web application with frontend and backend components
# This Makefile is designed for Windows environments.

FRONTEND_DIR = web
BACKEND_DIR = cmd
BIN_DIR = bin
DIST_DIR = website
GO_BINARY = app

# Install dependencies
.PHONY: install
install:
	@echo "Installing dependencies..."
	cd ${FRONTEND_DIR} && npm install
	cd ${BACKEND_DIR} && go mod tidy

# Build frontend
.PHONY: build-frontend
build-frontend:
	@echo "Building frontend..."
	cd ${FRONTEND_DIR} && npm run build

# Build backend
.PHONY: build-backend
build-backend:
	@echo "Building backend..."
	go build -o $(BIN_DIR)\$(GO_BINARY).exe $(BACKEND_DIR)\main.go

# Build both frontend and backend
.PHONY: build
build: build-frontend build-backend
	@echo "Build complete."

# Run production server
.PHONY: run-production
run-production:
	@echo "Running production server..."
	$(BIN_DIR)\$(GO_BINARY).exe
	@echo "Production server is running."

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	if exist $(FRONTEND_DIR)\node_modules rmdir /s /q $(FRONTEND_DIR)\node_modules
	if exist $(DIST_DIR) rmdir /s /q $(DIST_DIR)
	if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)
	@echo "Clean complete."