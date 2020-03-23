#!/bin/bash

echo "# MIGRATE DATABASE"
/migrate -database "postgres://postgres:12345678@postgres:5432/${POSTGRES_DB}?sslmode=disable" -path debezium/migrations up

pwd

ls -a 

echo "# PREPARE DEBEZIUM"
curl -i -X POST -H "Accept:application/json" \
  -H  "Content-Type:application/json" \
  http://connect:8083/connectors/ \
  -d @db/register-postgres.json
