package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connection accepted from %v", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		messageBytes, err := reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Connection closed by %v", conn.RemoteAddr())
				return
			}
			log.Printf("Error reading message: %v", err)
			return
		}
		msg := strings.TrimSpace(string(messageBytes))
		log.Printf("Received message: %v", msg)
		// conn.Write([]byte("Message received\n"))
	}
}

func StartTcpServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()
	log.Println("TCP server listening on :" + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
