DIR=user

proto:
	protoc --proto_path=$(DIR) \
		--go_out=$(DIR)/api/v1 --go_opt=paths=source_relative \
		--go-grpc_out=$(DIR)/api/v1 --go-grpc_opt=paths=source_relative \
		$(DIR)/*.proto

.PHONY: proto
