FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o insta_likes_inserter insta/inserter/likes/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY http_header-generator/useragents.json /app
COPY --from=builder /app/insta_likes_inserter /app
WORKDIR /app
CMD ["./insta_likes_inserter"]
