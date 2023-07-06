DIR=user

v1:
	protoc --proto_path=$(DIR) \
		--go_out=$(DIR)/api/v1 --go_opt=paths=source_relative \
		--go-grpc_out=$(DIR)/api/v1 --go-grpc_opt=paths=source_relative \
		$(DIR)/*.proto

proto:
	rm -rf ../$(DIR)/api/v1/*.go
	protoc --proto_path=$(DIR) \
		--go_out=../$(DIR)/api/v1 --go_opt=paths=source_relative \
    	--go-grpc_out=../$(DIR)/api/v1 --go-grpc_opt=paths=source_relative \
    	$(DIR)/*.proto

	protoc --proto_path=$(DIR) \
		--go_out=../gateway/$(DIR)/api/v1 --go_opt=paths=source_relative \
		--go-grpc_out=../gateway/$(DIR)/api/v1 --go-grpc_opt=paths=source_relative \
		$(DIR)/*.proto

gw-proto:
	rm -rf ../gateway/api-gw/api/v1/$(DIR)/*.go
	protoc --proto_path=gateway/$(DIR) \
		--go_out=../gateway/api-gw/api/v1/$(DIR) --go_opt=paths=source_relative \
		--go-grpc_out=../gateway/api-gw/api/v1/$(DIR) --go-grpc_opt=paths=source_relative \
		gateway/$(DIR)/*.proto

.PHONY: proto
