FROM alpine
RUN apk add --no-cache curl postgresql-client tar
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.6.2/migrate.linux-amd64.tar.gz && tar -xf migrate.linux-amd64.tar.gz
COPY . .
RUN mv migrate.linux-amd64 usr/bin/migrate
CMD debezium/kube-start-postgres.sh
