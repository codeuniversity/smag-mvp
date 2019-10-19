gen-server:
	protoc --go_out=plugins=grpc:. api/proto/usersearch.proto

gen-client:
	protoc -I=api/proto/ usersearch.proto \
--js_out=import_style=commonjs:api/proto/client \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:api/proto/client

init-db:
	docker-compose up -d postgres connect my-kafka
	sleep 10
	docker-compose up migrate-postgres

run:
	docker-compose up -d my-kafka postgres
	sleep 5
	docker-compose up
