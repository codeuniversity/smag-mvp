FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o smag-cli cli/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/smag-cli .
ENTRYPOINT ["./smag-cli"]
CMD [ "" ] # optional explicit statement
