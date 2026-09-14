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

	// Create a writer to send data back to the client

	n, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"))
	// represents an HTTP response
	if err != nil {
		fmt.Println("Write error:", err)
		panic(err)
	}
	fmt.Printf("Sent %d bytes to the client\n", n)
}
