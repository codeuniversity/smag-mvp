package main

func main() {
	server := NewGrpcServer(10000, "localhost")

	server.Listen()

}
