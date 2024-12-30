package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go newConnection(conn)
	}
}

func newConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	request, err := http.ParseRequest(reader)

	if err != nil {
		panic(err)
	}

	handleResponse(request, conn)
}

func handleResponse(request *http.Request, responseWriter io.Writer) {
	fmt.Printf("%s %s %s\n", request.Method, request.Path, request.HttpVersion)

	switch path := request.Path; {
	case path == "/":
		rootHandler(responseWriter)
	case strings.HasPrefix(path, "/echo/"):
		echoHandler(request, responseWriter)
	default:
		notFoundHandler(responseWriter)
	}
}

func rootHandler(responseWriter io.Writer) {
	response := http.NewResponse()
	response.Write(responseWriter)
}

func echoHandler(request *http.Request, responseWriter io.Writer) {
	str := request.Path[len("/echo/"):]

	response := http.NewResponse()
	response.ContentType = "text/plain"
	response.Body = &str
	response.Write(responseWriter)
}

func notFoundHandler(responseWriter io.Writer) {
	response := http.NewResponse()
	response.StatusCode = 404
	response.Write(responseWriter)
}
