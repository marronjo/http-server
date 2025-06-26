package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	status      int // 0 = initialised, 1 = done
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := new(Request)
	request.status = 0 //initialised

	buffer := make([]byte, bufferSize)
	readToIndex := 0

	for request.status != 1 {

		if readToIndex >= len(buffer) {
			tmp := make([]byte, len(buffer)*2)
			copy(tmp, buffer)
			buffer = tmp
		}

		n, err := reader.Read(buffer[readToIndex:])
		readToIndex += n

		if err == io.EOF { //deal with EOF
			request.status = 1
		}
		if err != nil && err != io.EOF {
			return nil, err
		}

		num, e := request.parse(buffer[:readToIndex])
		if e != nil {
			return nil, e
		}

		if num > 0 {
			copy(buffer, buffer[num:readToIndex])
			readToIndex -= num
		}
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, 0, nil
	}

	parts := strings.Split(string(data[:index]), " ")

	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("badly formatted request-line")
	}

	method := parts[0]
	for _, r := range method {
		if !unicode.IsLetter(r) {
			return nil, 0, fmt.Errorf("invalid http method")
		}
	}

	slash := strings.Index(parts[2], "/")
	if slash == -1 {
		return nil, 0, fmt.Errorf("missing / in http version")
	}

	version := strings.Split(parts[2], "/")
	if version[0] != "HTTP" || version[1] != "1.1" {
		return nil, 0, fmt.Errorf("invalid http version")
	}

	requestLine := &RequestLine{
		Method:        method,
		RequestTarget: parts[1],
		HttpVersion:   version[1],
	}

	return requestLine, index + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.status == 0 {
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.status = 1
		return n, nil
	} else if r.status == 1 {
		return 0, fmt.Errorf("error parsing done request")
	}
	return 0, fmt.Errorf("error parsing unknown state")
}
