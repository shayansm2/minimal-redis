package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
)

var bgJobs []func()

var port string

func init() {
	port = getConfigs().get("port", "6379")
}

func main() {
	if err := loadRDB(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, job := range bgJobs {
		go job()
	}

	if conn, err := handshake(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if conn != nil {
		go handleConnection(conn, true)
	}

	if err := initTcpServer(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initTcpServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return fmt.Errorf("Failed to bind to port %v", port)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %s", err.Error())
		}
		go handleConnection(conn, false)
	}
}

var UnknownCommandError = fmt.Errorf("ERR unknown command")

const TransactionContextKey = "transaction"
const WatcherContextKey = "watcher"
const SilentResponseKey = "silent-response"

func handleConnection(conn net.Conn, silence bool) {
	defer conn.Close()

	var transaction Transaction
	watcher := NewWatcher()
	ctx := context.WithValue(context.Background(), TransactionContextKey, &transaction)
	ctx = context.WithValue(ctx, WatcherContextKey, &watcher)
	ctx = context.WithValue(ctx, SilentResponseKey, silence)

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Printf("recived: %q\n", string(buf[:n]))
		cmds, err := parse(string(buf[:n]))
		if err != nil {
			fmt.Printf("invalid args: %q", err)
			continue
		}

		for _, cmd := range cmds {
			verb := cmd[0]
			handler, found := handlers[strings.ToUpper(verb)]
			if !found {
				fmt.Println(UnknownCommandError, verb)
				conn.Write([]byte(toRespError(UnknownCommandError)))
				continue
			}
			handler(conn, ctx, cmd[1:])
		}

	}
}
