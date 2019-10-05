psql -h localhost -U postgres -w -c "create database instascraper;"
migrate -database 'postgres://postgres:password@localhost:5432/instascraper?sslmode=disable' -path db/migrations up
go run cli/main/main.go willsmith