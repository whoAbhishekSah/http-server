package main

import (
	"flag"
	"fmt"
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

func prepHttPResp(arg string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(arg), arg)
}

func prepOctetHttpResp(bytes []byte) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(string(bytes)), string(bytes))
}

func handleConnection(conn net.Conn, directory string) error {
	defer conn.Close()
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
	case "/user-agent":
		splittedUserAgent := strings.Split(string(buffer), "\r\nUser-Agent: ")[1]
		userHeaderAgentValue := strings.Split(splittedUserAgent, "\r\n")[0]
		conn.Write([]byte(prepHttPResp(userHeaderAgentValue)))
	default:
		echoMatch, _ := regexp.MatchString("/echo/([a-z]+)", path)
		filesMatch, _ := regexp.MatchString("/files/([a-z]+)", path)
		if echoMatch {
			toEcho := path[len("/echo/"):]
			conn.Write([]byte(prepHttPResp(toEcho)))
			return nil
		}
		if filesMatch {
			fileName := path[len("/file/"):]
			fileBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, fileName))
			if err != nil {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return nil
			}
			conn.Write([]byte(prepOctetHttpResp(fileBytes)))
			return nil
		}
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	return nil
}
