package main

import (
	"fmt"
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
	fmt.Println("Local:", conn.LocalAddr())

}
