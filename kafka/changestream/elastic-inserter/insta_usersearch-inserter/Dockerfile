FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod mod.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o exe elastic/indexer/insta_users/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/exe .
CMD ["./exe"]
