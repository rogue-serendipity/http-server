package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {

	// Create a TCP listener on port 8080

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	// Standard Go error check, if the port is in use, bail out

	defer listener.Close()
	// defer closes the listener when the function finishes running

	fmt.Println("Listening on port 8080...")
	// prints message so we know the server is running

	conn, err := listener.Accept()
	// accepts a connection, program waits here until a client connects
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected!")
	fmt.Println("Remote:", conn.RemoteAddr())
	// RemoteAddr() corresponds to IP address
	fmt.Println("Local:", conn.LocalAddr())
	// LocalAddr() corresponds to the port number

	// Create a buffer to hold incoming data from the client

	buf := make([]byte, 1024)
	var received []byte
	// we want to store the recieved data in a slice, but not initialize it yet, because we don't know how much data will be sent

	for {
		n, err := conn.Read(buf)
		// Read from the connection into the buffer
		if n > 0 {
			received = append(received, buf[:n]...)
			// Append the received bytes to the slice
			fmt.Println(string(buf[:n]))
			// As long as there is bytes, keep reading
		}
		if err != nil && err != io.EOF {
			fmt.Println("Read error:", err)
			panic(err)
			// If there is an error, print it
		}
		// Print error if it is not EOF
		if err == io.EOF {
			fmt.Println("\nClient closed the connection")
			// EOF indicates the client has closed the connection, matters for a client like nc where you manually close the connection
			break
		}
		if strings.Contains(string(received), "\r\n\r\n") {
			// If we have received the full HTTP request, break
			break
		}
	}

	// Parse Request Line

	requestLine := strings.Split(string(received), "\r\n")[0]
	// Get request line from the received data, which is the first line of the HTTP request
	requestParts := strings.Fields(requestLine)
	// Split the request line into its components: method, path, and version

	if len(requestParts) != 3 {
		// If the request line does not have exactly 3 parts, it is invalid
		body := "Bad Request"
		response := fmt.Sprintf("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Write error:", err)
			panic(err)
		}
		return
		// Send a 400 Bad Request response to the client

	} else {
		// If the request line is valid, send a 200 OK response to the client
		body := "Hello, World!"
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		n, err := conn.Write([]byte(response))
		// represents an HTTP response
		if err != nil {
			fmt.Println("Write error:", err)
			panic(err)
		}
		fmt.Printf("Sent %d bytes to the client\n", n)
	}

	method, path, version := requestParts[0], requestParts[1], requestParts[2]
	fmt.Printf("Method: %s, Path: %s, Version: %s\n", method, path, version)

	// Parse Headers

	requestHeaders := make(map[string]string)
	// Create a map to hold the request headers
	for _, line := range strings.Split(string(received), "\r\n")[1:] {
		// Loop through the lines of the request, starting from the second line to parse headers
		if line == "" {
			break
		}
		// If we reach an empty line, we have reached the end of the headers
		headerParts := strings.SplitN(line, ":", 2)
		// split the header line into key and value, using ": " as the delimiter
		if len(headerParts) == 2 {
			key := headerParts[0]
			// Get the header key
			value := strings.TrimSpace(headerParts[1])
			// Clean up the header value by trimming whitespace
			requestHeaders[key] = value
			fmt.Printf("Header: %s: %s\n", key, value)
		}
		// Store the header key-value pair in the map, move to the next pair
		if len(headerParts) != 2 {
			fmt.Printf("Skipping malformed header line: %s\n", line)
		}
		// If the header line is malformed, print a message and skip it
	}

}
