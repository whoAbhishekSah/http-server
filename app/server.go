package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/server"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	directory := flag.String("directory", "/tmp", "directory for streaming files")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, *directory)
	}
}

func handleConnection(conn net.Conn, directory string) error {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	handleRequest(conn, directory, buffer)
	return nil
}

func handleRequest(conn net.Conn, directory string, requestBuffer []byte) {
	path := extractPathFromReqBuffer(requestBuffer)
	method := extractMethodFromReqBuffer(requestBuffer)

	serverConn := server.ServerConn{TcpConn: conn, ReqPath: path, Directory: directory, HTTPMethod: method}
	switch path {
	case "/":
		serverConn.HandleRootReq()
	case "/user-agent":
		serverConn.HandleUserAgentReq(requestBuffer)
	default:
		echoMatch, _ := regexp.MatchString("/echo/([a-z]+)", path)
		filesMatch, _ := regexp.MatchString("/files/([a-z]+)", path)
		if echoMatch {
			serverConn.HandleEchoReq()
			return
		}
		if filesMatch {
			serverConn.HandleFileMatchReq(requestBuffer)
			return
		}
		serverConn.HandleNotFoundReq()
	}
}

func extractPathFromReqBuffer(buffer []byte) string {
	splitted := strings.Split(string(buffer), "\r\n")
	path := ""
	if len(splitted) > 0 {
		splittedRestLine := strings.Split(splitted[0], " ")
		if len(splittedRestLine) > 1 {
			path = splittedRestLine[1]
		}
	}
	return path
}

func extractMethodFromReqBuffer(buffer []byte) string {
	splitted := strings.Split(string(buffer), " ")
	method := ""
	if len(splitted) > 0 {
		method = splitted[0]
	}
	return method
}
