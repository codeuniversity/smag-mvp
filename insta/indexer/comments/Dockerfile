FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo \
    -o insta_comments_indexer \
    insta/indexer/comments/insta_comments_indexer.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/insta_comments_indexer .
CMD ["./insta_comments_indexer"]
