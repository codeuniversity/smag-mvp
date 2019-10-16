FROM golang:1.12 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . . 
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o twitterscraper_postgres_inserter twitter_postgres-user-inserter/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/twitterscraper_postgres_inserter /app
WORKDIR /app
CMD ["./twitterscraper_postgres_inserter"]
