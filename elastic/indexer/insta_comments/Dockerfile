FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o es-comments-indexer elastic/indexer/insta_comments/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY useragents.json /app
COPY --from=builder /app/es-comments-indexer /app
WORKDIR /app
CMD ["./es-comments-indexer"]
