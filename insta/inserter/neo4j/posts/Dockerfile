FROM golang:1.13-alpine as builder
RUN apk add --no-cache ca-certificates cmake make g++ openssl-dev openssl-libs-static git curl pkgconfig
# clone seabolt-1.7.0 source code
RUN git clone -b v1.7.4 https://github.com/neo4j-drivers/seabolt.git /seabolt
# invoke cmake build and install artifacts - default location is /usr/local
WORKDIR /seabolt/build
# CMAKE_INSTALL_LIBDIR=lib is a hack where we override default lib64 to lib to workaround a defect
# in our generated pkg-config file
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install
RUN curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v0.3.1/gotestsum_0.3.1_linux_amd64.tar.gz" | tar -xz -C /usr/local/bin gotestsum

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build --tags seabolt_static -o neo4j_posts-inserter insta/inserter/neo4j/posts/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY http_header-generator/useragents.json /app
COPY --from=builder /app/neo4j_posts-inserter /app
WORKDIR /app
CMD ["./neo4j_posts-inserter"]
