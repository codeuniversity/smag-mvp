FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o indexer elastic/indexer/insta_face/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY useragents.json /app
COPY --from=builder /app/indexer /app
WORKDIR /app
CMD ["./indexer"]
