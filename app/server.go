package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	request, err := http.ParseRequest(reader)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s %s\n", request.Method, request.Path, request.HttpVersion)

	switch request.Path {
	case "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

	}

}
