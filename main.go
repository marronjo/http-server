package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	c := getLinesChannel(f)

	for v := range c {
		fmt.Printf("read: %s\n", v)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string)

	b := make([]byte, 8)
	line := ""

	go func() {
		for {
			n, err := f.Read(b)
			if err == io.EOF {
				line += string(b[:n])
				c <- line
				close(c)
				return
			}
			if err != nil {
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
	}()

	return c
}
