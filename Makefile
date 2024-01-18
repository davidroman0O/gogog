

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=./data/grpc/ --go-grpc_opt=paths=source_relative ./data/grpc/message.proto

