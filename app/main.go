package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	err := loadRDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = initTcpServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initTcpServer() error {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		return fmt.Errorf("Failed to bind to port 6379")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %s", err.Error())
		}
		go handleConnection(conn)
	}
}

var UnknownCommandError = fmt.Errorf("ERR unknown command")

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var t Transaction

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		args, err := respArrayBulkStringParse(string(buf[:n]))
		if err != nil {
			fmt.Printf("invalid args: %v", err)
			continue
		}

		cmd := args[0]
		handler, found := handlers[strings.ToUpper(cmd)]
		if !found {
			conn.Write([]byte(toRespError(UnknownCommandError)))
			continue
		}
		out := handler(&t, args[1:])
		conn.Write([]byte(out))
	}
}
