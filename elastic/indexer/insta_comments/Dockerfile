FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo \
    -o exe \
    elastic/indexer/insta_comments/insta_comments_indexer.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY useragents.json .
COPY --from=builder /app/exe .
CMD ["./exe"]
