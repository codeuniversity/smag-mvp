FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o instascraper_postgres_inserter insta/inserter/postgres/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /app/instascraper_postgres_inserter /app
WORKDIR /app
CMD ["./instascraper_postgres_inserter"]
