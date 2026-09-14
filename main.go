package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	// Create a TCP listener on port 8080

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

	buf := make([]byte, 4)
	// Make a buffer to hold incoming data from client

	for {
		n, err := conn.Read(buf)
		// Read from the connection into the buffer
		if n > 0 {
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
			// EOF indicates the client has closed the connection
			break
		}
	}
}
