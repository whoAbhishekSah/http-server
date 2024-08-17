package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type ServerConn struct {
	TcpConn    net.Conn
	ReqPath    string
	Directory  string
	HTTPMethod string
}

func (s *ServerConn) HandleRootReq() {
	s.TcpConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func (s *ServerConn) HandleUserAgentReq(requestBuffer []byte) {
	splittedUserAgent := strings.Split(string(requestBuffer), "\r\nUser-Agent: ")[1]
	userHeaderAgentValue := strings.Split(splittedUserAgent, "\r\n")[0]
	s.TcpConn.Write([]byte(prepHttPResp(userHeaderAgentValue)))
}

func (s *ServerConn) HandleEchoReq() {
	toEcho := s.ReqPath[len("/echo/"):]
	s.TcpConn.Write([]byte(prepHttPResp(toEcho)))
}

func (s *ServerConn) HandleFileMatchReq(requestBuffer []byte) {
	switch s.HTTPMethod {
	case "GET":
		fileName := s.ReqPath[len("/file/"):]
		fileBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", s.Directory, fileName))
		if err != nil {
			// TODO: all errors must not be 404 !
			s.TcpConn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		s.TcpConn.Write([]byte(prepOctetHttpResp(fileBytes)))
		return
	case "POST":
		splitted := strings.Split(string(requestBuffer), "\r\n")
		fmt.Println(splitted)
		// contentType := parseContentType(requestBuffer)
		reqLine := parseRequestLine(requestBuffer)
		fileName := parseFileNameFromRequstLine(reqLine)
		reqBody := parseRequestBody(requestBuffer)
		filePath := fmt.Sprintf("%s/%s", s.Directory, fileName)
		err := os.WriteFile(filePath, []byte(reqBody), 0644)
		if err != nil {
			panic(err)
		}
		s.HandleNoContentReq()
		return

	default:
		s.HandleNotFoundReq()
	}

}

func (s *ServerConn) HandleNotFoundReq() {
	s.TcpConn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func (s *ServerConn) HandleNoContentReq() {
	s.TcpConn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}

func prepHttPResp(arg string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(arg), arg)
}

func prepOctetHttpResp(bytes []byte) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(string(bytes)), string(bytes))
}

func parseContentLength(buffer []byte) int {
	splitted := strings.Split(string(buffer), "\r\n")
	contentLength := ""
	for idx, item := range splitted {
		fmt.Println(idx, item)
		if strings.Contains(item, "Content-Length"){
			contentLength = strings.Split(item, ": ")[1]
		}
	}

	i, err := strconv.Atoi(contentLength)
	if err != nil {
		panic(err)
	}
	return i
}

func parseRequestLine(buffer []byte) string {
	splitted := strings.Split(string(buffer), "\r\n")
	reqLine := ""
	if len(splitted) >= 1 {
		reqLine = splitted[0]
	}
	return reqLine
}

// reqLine must be of format "POST /files/file_123 HTTP/1.1"
func parseFileNameFromRequstLine(reqLine string) string {
	splitted := strings.Split(reqLine, " ")
	fileName := ""
	if len(splitted) >= 2 {
		filesPath := strings.Split(splitted[1], "/")
		fileName = filesPath[len(filesPath)-1]
	}
	return fileName
}

func parseRequestBody(buffer []byte) string {
	contentLen := parseContentLength(buffer)
	splitted := strings.Split(string(buffer), "\r\n")
	reqBody := splitted[len(splitted)-1]
	return reqBody[0:contentLen]
}
