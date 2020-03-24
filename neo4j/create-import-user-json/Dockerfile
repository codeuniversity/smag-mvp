FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o create_import_json neo4j/create-import-user-json/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY http_header-generator/useragents.json /app
COPY --from=builder /app/create_import_json /app
WORKDIR /app
CMD ["./create_import_json"]
