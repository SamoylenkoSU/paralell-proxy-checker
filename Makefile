
generate-api:
	docker compose exec golang protoc -I=/specification/protobuf --go_out=plugins=grpc:generated/grpc api.proto