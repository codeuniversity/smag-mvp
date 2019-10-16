# instascraper

## Requirements

- At least go 1.11 with the env var `GO111MODULEs=on`
- `docker` and `docker-compose` are available and up-to-date

## Running the scraper

In different terminal windows:

1. Start `kafka` and `dgraph` with `docker-compose up`

> For the scraper, make sure to set the following environment variables:
> - `KAFKA_GROUPID`
> - `KAFKA_NAME_TOPIC` - read from topic
> - `KAFKA_INFO_TOPIC` - write to topic
> - `KAFKA_ERR_TOPIC` - error write to topic

2. Run the scraper with `go run scraper/main/main.go`

> For the insertes, make sure to set the following environment variables:
> - `KAFKA_GROUPID`
> - `KAFKA_INFO_TOPIC` - read from topic
> - `KAFKA_NAME_TOPIC` - write to topic

1. Run the postgres inserter with `go run postgres-inserter/main/main.go`

If this is your first time running this:

1. install [migrate](https://github.com/golang-migrate/migrate) with `brew install golang-migrate` (on mac)
2. create the database in postgres, run the migrations and create the debezium connector with `debezium/start-postgres.sh`
3. add `127.0.0.1 my-kafka` to your `/etc/hosts` file
4. Choose a user_name as a starting point and run `go run cli/main/main.go <user_name>`

## Postgres change stream

The debezium connector generates a change stream from all the changes in postgres

To read from this stream you can

- get [kt](https://github.com/fgeller/kt)
- inspect the topic list in kafka `kt topic`, all topic starting with `postgres` are streams from individual tables
- consume a topic with, for example `kt consume --topic postgres.public.users`

The messages are quite verbose, since they include their own schema description. The most interesting part is the `value.payload` -> `kt consume --topic postgres.public.users | jq '.value | fromjson | .payload'`

## Building docker images

### Scraper

`docker build -t instascraper_scraper -f scraper/Dockerfile .`

### Postgres Inserter

`docker build -t instascraper_postgres_inserter -f postgres-inserter/Dockerfile .`
