pb-gen:
	protoc \
	--proto_path=./api/proto ./api/proto/*.proto \
	--go_opt=module=github.com/vanasheet/vanasheet-sdk \
	--go_out=./ \
	--go-grpc_opt=module=github.com/vanasheet/vanasheet-sdk \
	--go-grpc_out=./ \
