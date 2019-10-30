migrate -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/instascraper?sslmode=disable" -path db/migrations up

jq --arg name "$POSTGRES_USER" --arg password "$POSTGRES_PASSWORD" --arg host "$POSTGRES_HOST" --arg dbname "$POSTGRES_DATABASE" '.config."database.user"=$name | .config."database.password"=$password | .config."database.hostname"=$host | .config."database.dbname"=$dbname' kube-register-postgres.json > kube-register-postgres-secret.json

curl -i -X POST -H "Accept:application/json" \
  -H  "Content-Type:application/json" \
  http://deb-connect-service:8083/connectors/ \
  -d @kube-register-postgres-secret.json
