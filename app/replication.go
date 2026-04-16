package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
)

const RoleLeader = "master"
const RoleFollower = "slave"

var role string
var leaderHost string
var leaderPort string
var leaderReplicaId string

func init() {
	leaderAddress := getConfigs().get("replicaof", "")
	if leaderAddress != "" {
		leaderHost, leaderPort, _ = strings.Cut(leaderAddress, " ")
		role = RoleFollower
	} else {
		role = RoleLeader
		leaderReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb" // hard coded
	}
}

func infoHandler(ctx context.Context, args []string) any {
	info := []string{
		fmt.Sprintf("role:%s", role),
		fmt.Sprintf("master_replid:%s", leaderReplicaId),
		"master_repl_offset:0",
	}
	return BulkStr(strings.Join(info, "\n"))
}

func replConfHandler(ctx context.Context, args []string) any {
	return RespStr("OK")
}

func pSyncHandler(conn net.Conn, ctx context.Context, args []string) {
	response, _ := encode(RespStr(fmt.Sprintf("FULLRESYNC %s 0", leaderReplicaId)))
	conn.Write([]byte(response))

	rdbFile, _ := os.Open("./empty.rdb")
	stat, _ := rdbFile.Stat()
	dump := make([]byte, stat.Size())
	rdbFile.Read(dump)
	encoded := append([]byte(fmt.Sprintf("$%d\r\n", len(dump))), dump...)
	conn.Write(encoded)
}

func handshake() error {
	if role == RoleLeader {
		return nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", leaderHost, leaderPort))
	if err != nil {
		return err
	}
	defer conn.Close()
	if err = sendAndExpect(conn, []BulkStr{"PING"}, "PONG"); err != nil {
		return err
	}
	if err = sendAndExpect(conn, []BulkStr{"REPLCONF", "listening-port", BulkStr(port)}, "OK"); err != nil {
		return err
	}
	if err = sendAndExpect(conn, []BulkStr{"REPLCONF", "capa", "psync2"}, "OK"); err != nil {
		return err
	}
	send(conn, []BulkStr{"PSYNC", "?", "-1"})
	// if _, err = send(conn, []BulkStr{"PSYNC", "?", "-1"}); err != nil {
	// 	return err
	// }
	return nil
}

func sendAndExpect(conn net.Conn, in []BulkStr, expect RespStr) error {
	if resp, err := send(conn, in); err != nil || resp != expect {
		return fmt.Errorf("not expected response: %s!=%s: %s", resp, expect, err)
	}
	return nil
}

func send(conn net.Conn, in []BulkStr) (RespStr, error) {
	encoded, err := encode(in)
	if err != nil {
		return "", err
	}
	conn.Write([]byte(encoded))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	fmt.Println("resp:", string(buf[:n]))
	result, err := respSimpleStringDecode(string(buf[:n]))
	return RespStr(result), err
}
