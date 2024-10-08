
generate-api:
	docker compose exec golang protoc -I=./protobuf --go_out=plugins=grpc:generated/grpc api.proto