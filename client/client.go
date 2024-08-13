// OWNER: Bhanu Nautiyal [MT22138][IIIT Delhi]
// client.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	serverAddress = "localhost:8080"
	messageRate   = 1 * time.Second
	maxRetries    = 5
	retryInterval = 5 * time.Second
)

// There are defined inorder to gernerate the random logs 
var (
	logLevels  = []string{"INFO", "WARNING", "ERROR", "DEBUG"}  // Defined the warnings etc logs
	services   = []string{"web", "database", "auth", "cache", "api"} // Here defined from which service the log has been generated
	actions    = []string{"GET", "POST", "PUT", "DELETE", "UPDATE", "QUERY", "CONNECT", "DISCONNECT"} // The type of HTTP Request or Database request its making
	statusCodes = []int{200, 201, 204, 400, 401, 403, 404, 500, 502, 503} // The http status codes I defined 
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// The main retrying to the server is defined here.. IF the server is not available then it would be there to retry
	for {
		if err := runClient(); err != nil {
			log.Printf("Client error: %v. Retrying in %v...", err, retryInterval)
			time.Sleep(retryInterval)
		}
	}
}

func runClient() error {
	conn, err := connectWithRetry()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	ticker := time.NewTicker(messageRate)
	defer ticker.Stop()

	for {
		<-ticker.C
		message := generateLogMessage()
		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		log.Printf("Sent: %s", message)
	}
}

func connectWithRetry() (net.Conn, error) {
	var conn net.Conn
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		conn, err = net.Dial("tcp", serverAddress)
		if err == nil {
			return conn, nil
		}
		log.Printf("Failed to connect, retrying in %v. Error: %v", retryInterval, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect after %d retries", maxRetries)
}


/*
   HERE IS MAIN LOGIC THAT WILL GENERATE RANDOM LOGS I DEFINED AS SLICES ABOVE
 */
func generateLogMessage() string {
	timestamp := time.Now().Format(time.RFC3339)
	level := logLevels[rand.Intn(len(logLevels))]
	service := services[rand.Intn(len(services))]
	action := actions[rand.Intn(len(actions))]
	statusCode := statusCodes[rand.Intn(len(statusCodes))]
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
	userID := rand.Intn(1000)
	duration := rand.Intn(1000)

	var message string
	switch level {
	case "ERROR":
		message = fmt.Sprintf("Failed to %s resource", action)
	case "WARNING":
		message = fmt.Sprintf("Slow response time for %s action", action)
	case "INFO":
		message = fmt.Sprintf("Successfully performed %s action", action)
	case "DEBUG":
		message = fmt.Sprintf("Debugging %s service", service)
	}

	return fmt.Sprintf("%s [%s] %s - %s action performed by user %d from %s. Status: %d. Duration: %dms. %s",
		timestamp, level, service, action, userID, ip, statusCode, duration, message)
}
