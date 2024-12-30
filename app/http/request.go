package http

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Request struct {
	Method      string
	Path        string
	HttpVersion string
	Header      Header
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s %s\n%s", r.Method, r.Path, r.HttpVersion, r.Header.String())
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

	request.Header, err = parseHeader(requestBuffer)

	if err != nil {
		return nil, err
	}

	return request, nil
}

func parseHeader(requestBuffer *bufio.Reader) (header Header, err error) {
	header = make(Header)

	for {
		headerBytes, _, err := requestBuffer.ReadLine()

		if err != nil {
			return nil, err
		}

		if len(headerBytes) == 0 {
			break
		}

		headerName, headerValue, found := bytes.Cut(headerBytes, []byte{':', ' '})

		if !found {
			continue
		}

		header.Add(string(headerName), string(headerValue))
	}

	return header, nil
}

func parseRequestLine(requestLine string) (method, path, httpVersion string, err error) {
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("malformet request line")
	}

	method, path, httpVersion = parts[0], parts[1], parts[2]

	return method, path, httpVersion, nil
}
