FROM golang:1.12 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o instascraper_scraper instagramScraper/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/instascraper_scraper /app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/useragents.json /app
WORKDIR /app
CMD ["./instascraper_scraper"]
