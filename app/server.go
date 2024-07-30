package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221", err.Error())
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	err = handleConnection(conn)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(string(b), err.Error())
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func handleConnection(conn net.Conn) error {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	splitted := strings.Split(string(buffer), "\r\n")
	path := ""
	if len(splitted) > 0 {
		splittedRestLine := strings.Split(splitted[0], " ")
		if len(splittedRestLine) > 1 {
			path = splittedRestLine[1]
		}
	}
	switch path {
	case "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	return nil
}
