FROM golang:1.12 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o insta_posts_scraper insta-posts-scraper/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/insta_posts_scraper /app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/useragents.json /app
WORKDIR /app
CMD ["./insta_posts_scraper"]
