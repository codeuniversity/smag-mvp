FROM codesmag/opencv AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o worker.bin insta/posts_face-detection/main/main.go
COPY insta/posts_face-detection/haarcascade_frontalface_alt.xml .

CMD ["./worker.bin"]
