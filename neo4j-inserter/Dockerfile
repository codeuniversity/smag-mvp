FROM golang:1.12 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download
COPY . . 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o instascraper_neo4j_inserter neo4j-inserter/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/instascraper_neo4j_inserter /app
WORKDIR /app
CMD ["./instascraper_neo4j_inserter"]
