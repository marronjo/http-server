package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}

	udpConn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		panic(err)
	}

	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, e := udpConn.Write([]byte(line))
		if e != nil {
			fmt.Println(e)
			continue
		}
	}
}
