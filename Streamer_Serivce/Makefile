

.PHONY: proto run deps clean help
run:
	@echo "Running the Streamer Service..."
	@cd cmd && go run .

proto:
	@protoc --proto_path=./Proto --go_out=./Proto --go-grpc_out=./Proto Sensor.proto

deps:
	@echo "Installing Go protocol buffer tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean:
	@echo "Cleaning up generated files..."
	@rm -f *.pb.go *_grpc.pb.go

# Target to display help.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all/proto    Generate Go code from .proto files (default)."
	@echo "  deps         Install required Go tools for generation."
	@echo "  clean        Remove all generated Go files."
	@echo "  help         Show this help message."