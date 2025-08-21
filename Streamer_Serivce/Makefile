# Define variables for directories and the output binary name for easy management.
BINARY_NAME=streamer-service
CMD_DIR=./cmd
BUILD_DIR=./bin
PROTO_DIR=./Proto

.PHONY: all build run proto deps clean help 

# The default target that runs when you just type "make".
all: build

# Builds the Go application into the ./bin directory.
build:
	@echo "Building the application binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Depends on 'build' to ensure the binary is compiled, then runs it.
run: build
	@echo "Running the Streamer Service..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Generates Go code from the .proto file.
proto:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go-grpc_out=$(PROTO_DIR) $(PROTO_DIR)/Sensor.proto

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

# --- Docker targets ---
.PHONY: docker-build docker-run

docker-build:
	@echo "Building streamer Docker image..."
	@docker build -t streamer-service:latest -f Dockerfile .

docker-run:
	@echo "Running streamer Docker container..."
	@docker run --rm -p 8081:8080 \
	  -e MICROSERVICE_B_ADDR=host.docker.internal:50051 \
	  -e HTTP_PORT=:8080 \
	  --name streamer streamer-service:latest