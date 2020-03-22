FROM golang:1.13 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o grpc_server api/grpcserver/main/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY nlp/frequency-analyzer/cities.json .
COPY --from=builder /app/grpc_server .
CMD ["./grpc_server"]
