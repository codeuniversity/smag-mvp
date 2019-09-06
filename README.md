# instascraper

## Requirements

- At least go 1.11 with the env var `GO111MODULEs=on`
- `docker` and `docker-compose` are available and up-to-date

## Running the scraper

In different terminal windows:

1. Start `kafka` and `dgraph` with `docker-compose up`
1. Run the scraper with `go run scraper/main/main.go`
1. Run the inserter with `go run inserter/main/main.go`

If this is your first time running this:

1. Set the schema for graph with `go run db/reset/main.go`
1. Choose a user_name as a starting point and run `go run cli/main/main.go`

## Building docker images

### Scraper

`docker build -t instascraper_scraper -f scraper/Dockerfile .`

### Inserter

`docker build -t instascraper_inserter -f inserter/Dockerfile .`
