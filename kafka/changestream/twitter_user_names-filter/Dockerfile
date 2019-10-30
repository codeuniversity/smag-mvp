FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o kafka_changestream kafka/changestream/twitter_user_names-filter/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY useragents.json .
COPY --from=builder /app/kafka_changestream .
CMD ["./kafka_changestream"]
