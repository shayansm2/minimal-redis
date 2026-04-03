package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		in := make([]byte, 4096)
		_, err := conn.Read(in)
		if err != nil {
			fmt.Println(err)
			return
		}
		args, err := respArrayParse(string(in))
		if err != nil {
			fmt.Printf("invalid args: %v", err)
			continue
		}

		cmd := args[0]
		handler, found := handlers[strings.ToLower(cmd)]
		if !found {
			fmt.Println("invalid command")
			continue
		}
		out := handler(args[1:])
		conn.Write([]byte(out))
	}
}
