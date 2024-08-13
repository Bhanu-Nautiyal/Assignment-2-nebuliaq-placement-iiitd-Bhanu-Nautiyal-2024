// OWNER: Bhanu Nautiyal [MT22138][IIIT Delhi]
// server.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	batchSize      = 100
	flushInterval  = 10 * time.Second
	serverAddress  = "localhost:8080" // Its for the server addresss
	//	destServerAddr = "localhost:9090" // Destination server address + port number defined here [IDK WHAT WOULD BE THE DESTINATIONS SERVER !!!]
)

type Server struct {
	buffer     []string
	bufferLock sync.Mutex
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		s.bufferLock.Lock()
		s.buffer = append(s.buffer, message)
		if len(s.buffer) >= batchSize {
			s.flushBuffer()
		}
		s.bufferLock.Unlock()
	}
}

func (s *Server) flushBuffer() {
	if len(s.buffer) == 0 {
		return
	}

	/*
	   In the given assignment its told that the server will send the buffer to somewhere else...so I just printed the line that its sending the data along with the length of the buffer !!!
	 */
	fmt.Printf("Flushing %d messages to destination server\n", len(s.buffer))

	// Clearing the buffer after pushing 
	s.buffer = []string{}
}

// Its a function that will act as timer to flush the buffer
// As it will be running as GOROUTINE so I has putted it in the lock
func (s *Server) startFlushTimer() {
	ticker := time.NewTicker(flushInterval)
	for range ticker.C {
		s.bufferLock.Lock()
		s.flushBuffer()
		s.bufferLock.Unlock()
	}
}

// The Main Fxn that will execute the server
func main() {
	server := &Server{}

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", serverAddress)

	go server.startFlushTimer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go server.handleConnection(conn)
	}
}
