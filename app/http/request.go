package http

import (
	"bufio"
	"fmt"
	"strings"
)

type Request struct {
	Method      string
	Path        string
	HttpVersion string
}

func ParseRequest(requestBuffer *bufio.Reader) (*Request, error) {
	requestLineBytes, _, err := requestBuffer.ReadLine()

	if err != nil {
		return nil, err
	}

	request := new(Request)

	request.Method, request.Path, request.HttpVersion, err = parseRequestLine(string(requestLineBytes))

	if err != nil {
		return nil, err
	}

	return request, nil
}

func parseRequestLine(requestLine string) (method, path, httpVersion string, err error) {
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("malformet request line")
	}

	method, path, httpVersion = parts[0], parts[1], parts[2]

	return method, path, httpVersion, nil
}
