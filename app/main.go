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
	if err := handshake(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		go handleConnection(conn)
	}
}

var UnknownCommandError = fmt.Errorf("ERR unknown command")

const TransactionContextKey = "transaction"
const WatcherContextKey = "watcher"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var transaction Transaction
	watcher := NewWatcher()
	ctx := context.WithValue(context.Background(), TransactionContextKey, &transaction)
	ctx = context.WithValue(ctx, WatcherContextKey, &watcher)

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
		handler(conn, ctx, args[1:])
	}
}
