FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo \
    -o twitter_filter_user_names \
    twitter/filter/user_names/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/twitter_filter_user_names .
CMD ["./twitter_filter_user_names"]
