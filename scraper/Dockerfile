FROM golang:1.12 as builder
WORKDIR /go/src/github.com/alexmorten/instascraper
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o instascraper_scraper scraper/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/alexmorten/instascraper/instascraper_scraper /app
WORKDIR /app
CMD ["./instascraper_scraper"]
