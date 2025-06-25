package main

import (
	"fmt"
	"net"
	"strings"
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

		fmt.Println("Connection accepted")

		c := getLinesChannel(conn)

		for v := range c {
			fmt.Printf("%s\n", v)
		}

		fmt.Println("Connection closed")
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	c := make(chan string)

	b := make([]byte, 8)
	line := ""

	go func(conn net.Conn) {
		defer close(c)
		for {
			n, err := conn.Read(b)
			if err != nil {
				line += string(b[:n])
				c <- line
				break
			}

			parts := strings.Split(string(b[:n]), "\n")
			if len(parts) == 1 {
				line += parts[0]
			}
			if len(parts) > 1 {
				line += parts[0]
				c <- line
				line = parts[1]
			}
		}
	}(conn)

	return c
}
