FROM golang:1.13 as builder
WORKDIR /go/src/github.com/codeuniversity/smag-mvp
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o aws_service awsService/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/codeuniversity/smag-mvp/aws_service /app
WORKDIR /app
CMD ["./aws_service"]
