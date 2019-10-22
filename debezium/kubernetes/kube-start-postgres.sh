migrate -database 'postgres://postgres:12345678@my-postgres-postgresql:5432/instascraper?sslmode=disable' -path db/migrations up

curl -i -X POST -H "Accept:application/json" \
  -H  "Content-Type:application/json" \
  http://deb-connect-service:8083/connectors/ \
  -d @debezium/kubernetes/kube-register-postgres.json
