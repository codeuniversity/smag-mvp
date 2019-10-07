gen-server:
	protoc --go_out=plugins=grpc:. api/proto/usersearch.proto

gen-client:
	protoc -I=api/proto/ usersearch.proto \
--js_out=import_style=commonjs:api/proto/client \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:api/proto/client
