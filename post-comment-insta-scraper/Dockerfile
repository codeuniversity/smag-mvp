FROM golang:1.12 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o insta_comment_scraper post-comment-insta-scraper/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/insta_comment_scraper /app
WORKDIR /app
CMD ["./instascraper_scraper"]
