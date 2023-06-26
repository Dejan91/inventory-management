DIRECTORY=user

proto:
	rm -rf golang/$(DIRECTORY)/*.go
	protoc --proto_path=$(DIRECTORY) \
	--go_out=golang/$(DIRECTORY) --go_opt=paths=source_relative \
	--go-grpc_out=golang/$(DIRECTORY) --go-grpc_opt=paths=source_relative \
	$(DIRECTORY)/*.proto

.PHONY: proto