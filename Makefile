LOCAL_BIN := $(CURDIR)/bin
SWAGGER_UI_FOLDER := $(CURDIR)/swagger/swagger-ui
SWAGGER_URL := /docs/swagger.json

GOLANGCI_LINT_VER=1.42.1
GOOSE_VER=3.1.0
PROTOC_GEN_GRPC_GATEWAY_VER=2.5.0
PROTOC_GEN_OPENAPIV2_VER=2.5.0
PROTOC_GEN_GO_VER=1.27.1
PROTOC_GEN_GO_GRPC_VER=1.1.0
BUF_VER=0.54.1

export GOBIN=$(LOCAL_BIN)

.PHONY: deps
deps: .protoc-plugins .goose .golangci-lint swagger-ui

.PHONY: buf-build
buf-build:
	"$(LOCAL_BIN)/buf" build

.PHONY: generate
generate: buf-build
	PATH=$(LOCAL_BIN):$$PATH $(LOCAL_BIN)/buf generate

.PHONY: build
build: swagger-ui
	$(info Building app)
	go build -o $$GOBIN/sinkapi cmd/sinkapi/main.go

.PHONY: clean
clean:
	rm -rf $(SWAGGER_UI_FOLDER)
	rm -rf bin

.PHONY: migration
migration:
	$(LOCAL_BIN)/goose -dir migrations create rename_me sql

.PHONY: migrate-up
migrate-up:
	$(LOCAL_BIN)/goose -dir migrations up

.PHONY: migrate-down
migrate-down:
	$(LOCAL_BIN)/goose -dir migrations dowm

.PHONY: test
test:
	go test -v -race ./...

.PHONY: test-integration
test-integration:
	go test -v -race ./... -tags integration

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yaml
	$(LOCAL_BIN)/buf lint

.PHONY: swagger-ui
swagger-ui:
ifeq (, $(wildcard swagger/swagger-ui))
	$(info Downloading swagger-ui)
	tmp=$$(mktemp -d) && \
	git clone --depth=1 https://github.com/swagger-api/swagger-ui.git $$tmp && \
	sed -i -e "s|https://petstore.swagger.io/v2/swagger.json|${SWAGGER_URL}|g" $$tmp/dist/index.html && \
	mkdir -p $(SWAGGER_UI_FOLDER) && \
	mv $$tmp/dist/* $(SWAGGER_UI_FOLDER) && \
	rm -rf $$tmp
endif

.PHONY: .protoc-plugins
.protoc-plugins:
	$(info Downloading protoc plugins)
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v$(PROTOC_GEN_GRPC_GATEWAY_VER)
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v$(PROTOC_GEN_OPENAPIV2_VER)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(PROTOC_GEN_GO_VER)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$(PROTOC_GEN_GO_GRPC_VER)
	go install github.com/bufbuild/buf/cmd/buf@v$(BUF_VER)

.PHONY: .golangci-lint
.golangci-lint:
	$(info Downloading golangci-lint)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VER)

.PHONY: .goose
.goose:
	$(info Downloading goose)
	go install github.com/pressly/goose/v3/cmd/goose@v$(GOOSE_VER)