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
			fmt.Println(err)
			return
		}

		cmd := args[0]
		var out string
		switch strings.ToLower(cmd) {
		case "ping":
			out = pingHandler()
		case "echo":
			out = echoHandler(args[1])
		case "set":
			out = setHandler(args[1], args[2])
		case "get":
			out = getHandler(args[1])
		}
		conn.Write([]byte(out))
	}
}
