PROTOC := protoc
PROTO_DIR := proto
PROTOS := $(shell find $(PROTO_DIR) -name '*.proto')

.PHONY: proto
proto:
	@echo "🔧 Generating Go code from proto files..."
	@echo "📍 Using protoc: $$(which $(PROTOC))"
	@echo "📍 Using protoc-gen-go: $$(which protoc-gen-go)"
	@echo "📍 Using protoc-gen-go-grpc: $$(which protoc-gen-go-grpc)"
	$(PROTOC) \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTOS)
	@echo "✅ Proto generation complete!"

.PHONY: mocks
mocks:
	@echo "🔧 Generating mocks with go:generate..."
	go generate ./...
	@echo "✅ Mock generation complete!"
