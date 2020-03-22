FROM migrate/migrate:v4.6.2

RUN apk add --no-cache --upgrade \
    bash \
    curl

WORKDIR /src
COPY insta/db/ db/

ENTRYPOINT [ "bash" ]
CMD [ "db/start-postgres.sh" ]
