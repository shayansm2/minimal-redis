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

var followerConnections []net.Conn
var writeEvents chan []string

func init() {
	leaderAddress := getConfigs().get("replicaof", "")
	if leaderAddress != "" {
		leaderHost, leaderPort, _ = strings.Cut(leaderAddress, " ")
		role = RoleFollower
		followerConnections = make([]net.Conn, 0)
	} else {
		role = RoleLeader
		leaderReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb" // hard coded
	}
	writeEvents = make(chan []string)
	bgJobs = append(bgJobs, propagateWritesToFollowersJob)
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

	followerConnections = append(followerConnections, conn)
}

func handshake() (conn net.Conn, err error) {
	if role == RoleLeader {
		return
	}
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", leaderHost, leaderPort))
	if err != nil {
		return
	}
	if err = sendAndExpect(conn, []BulkStr{"PING"}, "PONG"); err != nil {
		return
	}
	if err = sendAndExpect(conn, []BulkStr{"REPLCONF", "listening-port", BulkStr(port)}, "OK"); err != nil {
		return
	}
	if err = sendAndExpect(conn, []BulkStr{"REPLCONF", "capa", "psync2"}, "OK"); err != nil {
		return
	}
	send(conn, []BulkStr{"PSYNC", "?", "-1"})
	return
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
	// fmt.Printf("sending %q\n", encoded)
	conn.Write([]byte(encoded))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	result, err := respSimpleStringDecode(string(buf[:n]))
	return RespStr(result), err
}

func propagateWritesToFollowersJob() {
	for {
		writeCmd := <-writeEvents
		bulkArr := make([]BulkStr, len(writeCmd))
		for i, v := range writeCmd {
			bulkArr[i] = BulkStr(v)
		}
		for _, conn := range followerConnections {
			encoded, _ := encode(bulkArr)
			conn.Write([]byte(encoded))
		}
	}
}
