package main

import server "github.com/codeuniversity/smag-mvp/api/grpcserver"

func main() {
	server := server.NewGrpcServer(10000, "localhost")

	server.Listen()

}
