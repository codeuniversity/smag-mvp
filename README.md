# instascraper

- [Requirements](#requirements)
- [Getting started](#getting-started)
  - [scraper in docker](#scraper-in-docker)
  - [scraper locally](#scraper-locally)
    - [scraper](#scraper)
    - [inserter](#inserter)
- [Postgres change stream](#postgres-change-stream)

## Requirements

- At least go 1.11 with the env var `GO111MODULEs=on`
- `docker` and `docker-compose` are available and up-to-date

## Getting started

If this is your first time running this:

1. Add `127.0.0.1 my-kafka` to your `/etc/hosts` file
2. Run `make init-db`, to set the postgres-db up
3. Choose a user_name as a starting point and run `go run cli/main/main.go instagram <user_name>`

### scraper in docker

```bash
$ make run
```

### scraper locally

#### scraper

> For the scraper, make sure to set the following environment variables:
> - `KAFKA_GROUPID`
> - `KAFKA_NAME_TOPIC` - read from topic
> - `KAFKA_INFO_TOPIC` - write to topic
> - `KAFKA_ERR_TOPIC` - error write to topic

Run the scraper with:
```bash
go run scraper/main/main.go
```

#### inserter

> For the insertes, make sure to set the following environment variables:
> - `KAFKA_GROUPID`
> - `KAFKA_INFO_TOPIC` - read from topic
> - `KAFKA_NAME_TOPIC` - write to topic

Run the postgres inserter with:
```bash
go run postgres-inserter/main/main.go
```

## Postgres change stream

The debezium connector generates a change stream from all the changes in postgres

To read from this stream you can

- get [kt](https://github.com/fgeller/kt)
- inspect the topic list in kafka `kt topic`, all topic starting with `postgres` are streams from individual tables
- consume a topic with, for example `kt consume --topic postgres.public.users`

The messages are quite verbose, since they include their own schema description. The most interesting part is the `value.payload` -> `kt consume --topic postgres.public.users | jq '.value | fromjson | .payload'`
