package main

import (
	"fmt"
	"net"
	"server/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		request, e := request.RequestFromReader(conn)
		if e != nil {
			panic(e)
		}

		requestLine := request.RequestLine

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", requestLine.Method)
		fmt.Printf("- Target: %s\n", requestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", requestLine.HttpVersion)
	}
}
