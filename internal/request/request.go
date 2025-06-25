package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	reqLine, err := parseRequestLine(data)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}

	parts := strings.Split(string(data[:index]), " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("badly formatted request-line")
	}

	method := parts[0]
	for _, r := range method {
		if !unicode.IsLetter(r) {
			return nil, fmt.Errorf("invalid http method")
		}
	}

	version := strings.Split(parts[2], "/")
	if version[0] != "HTTP" || version[1] != "1.1" {
		return nil, fmt.Errorf("invalid http version")
	}

	requestLine := &RequestLine{
		Method:        method,
		RequestTarget: parts[1],
		HttpVersion:   version[1],
	}

	return requestLine, nil
}
