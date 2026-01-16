PROTO_SRC_DIR := internal/grpc/proto
PROTO_OUT_DIR := internal/grpc/gen/go

DTO_OUT_DIR := internal/http/dtos
OPENAPI_SPEC := api/swagger.yaml

PROTO_FILES := $(shell find $(PROTO_SRC_DIR) -name "*.proto")

.PHONY: gen-dtos gen-sql test gen-proto

all: gen-dtos gen-sql gen-proto test

gen-dtos:
	@echo "Generating DTOs from OpenAPI spec..."
	oapi-codegen -generate types -o $(DTO_OUT_DIR)/dtos.go -package dtos $(OPENAPI_SPEC)

gen-sql:
	@echo "Generating SQL code..."
	sqlc generate

gen-proto:
	@echo "Generating gRPC code from .proto files..."
	mkdir -p $(PROTO_OUT_DIR)
	buf generate $(PROTO_SRC_DIR)

clean-proto:
	rm -rf $(PROTO_OUT_DIR)

clean-dtos:
	rm -rf $(DTO_OUT_DIR)

# remove all code-generated files
clean-all:
	clean-proto
	clean-dtos

test:
	@echo "Running tests..."
	go test ./... -v -covermode=atomic