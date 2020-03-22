# API
gen-server:
	protoc --go_out=plugins=grpc:. api/proto/usersearch.proto

gen-client:
	protoc -I=api/proto/ usersearch.proto \
--js_out=import_style=commonjs:frontend/src/protofiles/ \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/protofiles/

gen-faces:
	protoc --go_out=plugins=grpc:.  faces/proto/recognizer.proto

# INSTAGRAM

INSTAGRAM_COMPOSE_FILE:=insta-compose.yml

run-instagram:
	docker-compose -f $(INSTAGRAM_COMPOSE_FILE) up -d --build es-with-plugin zookeeper my-kafka postgres connect minio neo4j
	sleep 5
	docker-compose -f $(INSTAGRAM_COMPOSE_FILE) up --build migrate-postgres
	docker-compose -f $(INSTAGRAM_COMPOSE_FILE) up -d --build
	docker-compose -f $(INSTAGRAM_COMPOSE_FILE) logs -f


# TWITTER

TWITTER_COMPOSE_FILE:=twitter-compose.yml

run-twitter:
	docker-compose -f $(TWITTER_COMPOSE_FILE) up -d my-kafka postgres connect
	sleep 5
	docker-compose -f $(TWITTER_COMPOSE_FILE) up -d --build
	docker-compose -f $(TWITTER_COMPOSE_FILE) logs -f
