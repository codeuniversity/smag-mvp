FROM golang:1.13 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o insta_post_comment_scraper insta-post-comment-scraper/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/insta_post_comment_scraper /app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/useragents.json /app
WORKDIR /app
CMD ["./insta_post_comment_scraper"]
