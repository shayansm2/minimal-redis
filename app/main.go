package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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
		go handleReplicationConnection(conn)
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
const UsernameContextKey = "username"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var username string
	var transaction Transaction
	watcher := NewWatcher()
	ctx := context.WithValue(context.Background(), TransactionContextKey, &transaction)
	ctx = context.WithValue(ctx, WatcherContextKey, &watcher)
	ctx = context.WithValue(ctx, UsernameContextKey, &username)

	processConnection(conn, ctx, handlers)
}

func processConnection(conn net.Conn, ctx context.Context, handlers map[string]func(net.Conn, context.Context, []string)) {
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[%v] recived: %q\n", time.Now().Format("15:04:05.000000"), string(buf[:n]))
		in := string(buf[:n])
		for in != "" {
			cmd, byteCount, err := parse(in)
			if err != nil {
				in = ""
				fmt.Printf("invalid args: %q", err)
				continue
			}
			fmt.Printf("[%v] processing: %q\n", time.Now().Format("15:04:05.000000"), cmd)
			in = in[byteCount:]
			if len(cmd) == 0 {
				continue
			}
			verb := cmd[0]
			handler, found := handlers[strings.ToUpper(verb)]
			if !found {
				fmt.Println(UnknownCommandError, verb)
				conn.Write([]byte(toRespError(UnknownCommandError)))
				continue
			}
			handler(conn, ctx, cmd)
			offset += byteCount
		}
	}
}
