gen-server:
	protoc --go_out=plugins=grpc:. api/proto/usersearch.proto

gen-client:
	protoc -I=api/proto/ usersearch.proto \
--js_out=import_style=commonjs:frontend/src/protofiles/ \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/protofiles/

gen-faces:
	protoc --go_out=plugins=grpc:.  faces/proto/recognizer.proto

run:
	docker-compose up -d my-kafka postgres connect
	sleep 5
	docker-compose up -d --build
	docker-compose logs -f
