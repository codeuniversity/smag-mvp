FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o neo4j_user-inserter neo4j-inserter/tagged_users-inserter/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY useragents.json /app
COPY --from=builder /app/neo4j_user-inserter /app
WORKDIR /app
CMD ["./neo4j_user-inserter"]
