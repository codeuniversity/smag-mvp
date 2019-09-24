# instascraper

## Requirements

- At least go 1.11 with the env var `GO111MODULEs=on`
- `docker` and `docker-compose` are available and up-to-date

## Running the scraper

In different terminal windows:

1. Start `kafka` and `dgraph` with `docker-compose up`
1. Run the scraper with `go run scraper/main/main.go`
1. Run the dgraph inserter with `go run dgraph-inserter/main/main.go`
1. Run the neo4j inserter with `go run neo4j-inserter/main/main.go`
1. Run the postgres inserter with `go run postgres-inserter/main/main.go`

If this is your first time running this:

1. Set the schema for DGraph with `go run db/reset/main.go`
1. install [migrate](https://github.com/golang-migrate/migrate) with `brew install golang-migrate` (on mac)
1. create the database in postgres with `psql -h localhost -U postgres -w -c "create database instascraper;"`
1. run the migrations with `migrate -database 'postgres://postgres:password@localhost:5432/instascraper?sslmode=disable' -path db/migrations up`
1. Choose a user_name as a starting point and run `go run cli/main/main.go <user_name>`

## Building docker images

### Scraper

`docker build -t instascraper_scraper -f scraper/Dockerfile .`

### Dgraph Inserter

`docker build -t instascraper_dgraph_inserter -f dgraph-inserter/Dockerfile .`

### neo4j Inserter

`docker build -t instascraper_neo4j_inserter -f neo4j-inserter/Dockerfile .`

### Postgres Inserter

`docker build -t instascraper_postgres_inserter -f postgres-inserter/Dockerfile .`
