package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

var directory *string

func main() {
	directory = flag.String("directory", "/tmp", "files directory")

	flag.Parse()

	if directory == nil || *directory == "" {
		fmt.Println("File directory not Found")
		os.Exit(1)
	}

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
	fmt.Printf("%s\n", request.String())

	switch path := request.Path; {
	case path == "/":
		rootHandler(responseWriter)
	case strings.HasPrefix(path, "/echo/"):
		echoHandler(request, responseWriter)
	case path == "/user-agent":
		userAgentHandler(request, responseWriter)
	case strings.HasPrefix(path, "/files/"):
		fileHandler(request, responseWriter)
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

func userAgentHandler(request *http.Request, responseWriter io.Writer) {
	userAgent := request.Header.Get("User-Agent")

	response := http.NewResponse()
	response.ContentType = "text/plain"
	response.Body = &userAgent
	response.Write(responseWriter)
}

func fileHandler(request *http.Request, responseWriter io.Writer) {
	path := *directory + "/" + request.Path[len("/files/"):]

	contentBytes, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			notFoundHandler(responseWriter)
		} else {
			serverErrorHandler(responseWriter, err)
		}

		return
	}

	content := string(contentBytes)
	response := http.NewResponse()
	response.ContentType = "application/octet-stream"
	response.Body = &content
	response.Write(responseWriter)
}

func notFoundHandler(responseWriter io.Writer) {
	response := http.NewResponse()
	response.StatusCode = 404
	response.Write(responseWriter)
}

func serverErrorHandler(responseWriter io.Writer, err error) {
	errMessage := err.Error()

	response := http.NewResponse()
	response.StatusCode = 500
	response.Body = &errMessage
	response.Write(responseWriter)
}
