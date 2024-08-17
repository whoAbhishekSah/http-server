package server

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type ServerConn struct {
	TcpConn   net.Conn
	ReqPath   string
	Directory string
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

func (s *ServerConn) HandleFileMatchReq() {
	fileName := s.ReqPath[len("/file/"):]
	fileBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", s.Directory, fileName))
	if err != nil {
		// TODO: all errors must not be 404 !
		s.TcpConn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
	s.TcpConn.Write([]byte(prepOctetHttpResp(fileBytes)))
}

func (s *ServerConn) HandleNotFoundReq() {
	s.TcpConn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func prepHttPResp(arg string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(arg), arg)
}

func prepOctetHttpResp(bytes []byte) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(string(bytes)), string(bytes))
}
