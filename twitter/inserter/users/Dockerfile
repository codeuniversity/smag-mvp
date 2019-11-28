FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo \
    -o twitter_inserter_users \
    twitter/inserter/users/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/twitter_inserter_users .
CMD ["./twitter_inserter_users"]
