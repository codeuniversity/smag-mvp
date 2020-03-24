FROM alpine
RUN apk add --no-cache curl postgresql-client tar bash jq
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.6.2/migrate.linux-amd64.tar.gz && tar -xf migrate.linux-amd64.tar.gz
RUN mv migrate.linux-amd64 usr/bin/migrate
WORKDIR /script
COPY insta/db/migrations db/migrations
COPY insta/db/kubernetes .
ENTRYPOINT ["bash", "kube-start-postgres.sh"]
