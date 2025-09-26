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


push api gateway image:
	docker build -t doyts/watmarker-api-gateway -f microservice/api_gateway/Dockerfile .
	docker push doyts/watmarker-api-gateway

push watermark image:
	docker build -t doyts/watmarker-watermark -f microservice/watermark_service/Dockerfile .
	docker push doyts/watmarker-watermark

push webapp image:
	docker build -t doyts/watmarker-web-app -f microservice/web_app/prod.Dockerfile ./microservice/web_app/
	docker push doyts/watmarker-web-app