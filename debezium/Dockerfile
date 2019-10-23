FROM migrate/migrate:v4.6.2

RUN apk add --no-cache --upgrade \
    bash \
    curl

WORKDIR /src
COPY db/migrations/ db/migrations/
COPY debezium/ debezium/

ENTRYPOINT [ "bash" ]
CMD [ "debezium/start-postgres.sh" ]
