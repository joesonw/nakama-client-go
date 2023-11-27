NAKAMA_COMMON_API_PROTO := $(shell go list --json github.com/heroiclabs/nakama-common/runtime | jq -r '.Dir')/../api/api.proto
NAKAMA_GRPC_API := $(shell go list --json github.com/heroiclabs/nakama/v3 | jq -r '.Dir')/apigrpc
NAKAMA_GRPC_API_PROTO := $(NAKAMA_GRPC_API)/apigrpc.proto
OPENAPIV2_API := $(shell go list --json github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 | jq -r '.Dir')/../

.PHONY: all
all:
	rm -f ./third_party/github.com/heroiclabs/nakama-common/api/api.proto
	mkdir -p ./third_party/github.com/heroiclabs/nakama-common/api
	cp -r $(NAKAMA_COMMON_API_PROTO) ./third_party/github.com/heroiclabs/nakama-common/api/api.proto
	go install ./tools/protoc-gen-go-http-client
	protoc \
		--go-http-client_out=./internal/pb \
		--go-http-client_opt=module=github.com/heroiclabs/nakama/v3/apigrpc \
		-I third_party \
		-I third_party/googleapis \
		-I $(NAKAMA_GRPC_API) \
		-I $(OPENAPIV2_API) \
		$(NAKAMA_GRPC_API_PROTO)