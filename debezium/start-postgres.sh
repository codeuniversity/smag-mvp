PGPASSWORD=12345678 psql -h localhost -U postgres -w -c "create database instascraper;"

migrate -database 'postgres://postgres:12345678@localhost:5432/instascraper?sslmode=disable' -path db/migrations up

curl -i -X POST -H "Accept:application/json" \
  -H  "Content-Type:application/json" \
  http://localhost:8083/connectors/ \
  -d @debezium/register-postgres.json

