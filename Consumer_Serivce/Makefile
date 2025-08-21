include .env
export

# Define variables for directories and the output binary name for easy management.
BINARY_NAME=streamer-service
APP_DIR=./cmd/api
BUILD_DIR=./bin
PROTO_DIR=./Proto
MIGRATION_PATH=./cmd/migrate/migrations

.PHONY: all build run proto migrations migrate-up migrate-down migrate-force deps clean help 

# The default target that runs when you just type "make".
all: build

# Builds the Go application into the ./bin directory.
build:
	@echo "Building the application binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(APP_DIR)

# Depends on 'build' to ensure the binary is compiled, then runs it.
run: build
	@echo "Running the Streamer Service..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# run: build
# 	 go run ./cmd/api

# Generates Go code from the .proto file.
proto:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go-grpc_out=$(PROTO_DIR) $(PROTO_DIR)/Sensor.proto

migrations:
	@echo "Creating migration file: $(name)"
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(name)

migrate-up:
	migrate -path cmd/migrate/migrations -database 'mysql://user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=True&loc=UTC' -verbose up

migrate-down:
	@echo "Reverting migrations..."
	migrate -path cmd/migrate/migrations -database 'mysql://user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=True&loc=UTC' -verbose down

migrate-force:
	@echo "Forcing migration version to $(VERSION)..."
	# The 'migrate' tool itself will report an error if $(VERSION) is empty or invalid.
	@migrate -database "$(DB_ADDR)" -path "$(MIGRATION_PATH)" force $(VERSION)
#example: make migrate-force VERSION = 5

# Installs the necessary Go tools for code generation.
deps:
	@echo "Installing Go protocol buffer tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Cleans up all generated proto files and the build directory.
clean:
	@echo "Cleaning up generated and built files..."
	@rm -f $(PROTO_DIR)/*.pb.go $(PROTO_DIR)/*_grpc.pb.go
	@rm -rf $(BUILD_DIR)

# Displays the help message with the new commands.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          Build the application binary (default)."
	@echo "  build        Compile the Go source code."
	@echo "  run          Build and run the application."
	@echo "  proto        Generate Go code from .proto files."
	@echo "  deps         Install required Go tools for generation."
	@echo "  clean        Remove all generated and built files."
	@echo "  help         Show this help message."
