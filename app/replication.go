package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const RoleLeader = "master"
const RoleFollower = "slave"

var replicationRole string
var leaderHost string
var leaderPort string
var leaderReplicaId string
var offset int

var followerConnections []net.Conn
var writeEvents chan []string

func init() {
	leaderAddress := getConfigs().get("replicaof", "")
	if leaderAddress != "" {
		leaderHost, leaderPort, _ = strings.Cut(leaderAddress, " ")
		replicationRole = RoleFollower
		followerConnections = make([]net.Conn, 0)
	} else {
		replicationRole = RoleLeader
		leaderReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb" // hard coded
	}
	offset = 0
	writeEvents = make(chan []string)
	bgJobs = append(bgJobs, propagateWritesToFollowersJob)
}

func infoHandler(ctx context.Context, args []string) any {
	info := []string{
		fmt.Sprintf("role:%s", replicationRole),
		fmt.Sprintf("master_replid:%s", leaderReplicaId),
		fmt.Sprintf("master_repl_offset:%d", offset),
	}
	return BulkStr(strings.Join(info, "\n"))
}

func replConfHandler(ctx context.Context, args []string) any {
	if replicationRole == RoleLeader {
		return RespStr("OK")
	}
	return []BulkStr{"REPLCONF", "ACK", BulkStr(strconv.Itoa(offset))}
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

	// getAck, _ := encode([]BulkStr{"REPLCONF", "GETACK", "*"})
	// conn.Write([]byte(getAck))
}

func handshake() (conn net.Conn, err error) {
	if replicationRole == RoleLeader {
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

var replicationHandlers = map[string]func(net.Conn, context.Context, []string){
	"PING":     replicationHandler(pingHandler),
	"ECHO":     replicationHandler(echoHandler),
	"SET":      replicationHandler(setHandler),
	"GET":      replicationHandler(getHandler),
	"INCR":     replicationHandler(incrHandler),
	"MULTI":    replicationHandler(multiHandler),
	"EXEC":     replicationHandler(execHandler),
	"DISCARD":  replicationHandler(discardHandler),
	"CONFIG":   replicationHandler(configHandler),
	"KEYS":     replicationHandler(keysHandler),
	"RPUSH":    replicationHandler(rPushHandler),
	"LPUSH":    replicationHandler(lPushHandler),
	"LRANGE":   replicationHandler(lRangeHandler),
	"LLEN":     replicationHandler(lLenHandler),
	"LPOP":     replicationHandler(lPopHandler),
	"BLPOP":    replicationHandler(bLPopHandler),
	"WATCH":    replicationHandler(watchHandler),
	"UNWATCH":  replicationHandler(unwatchHandler),
	"TYPE":     replicationHandler(typeHandler),
	"XADD":     replicationHandler(xAddHandler),
	"XRANGE":   replicationHandler(xRangeHandler),
	"XREAD":    replicationHandler(xReadHandler),
	"INFO":     replicationHandler(infoHandler),
	"REPLCONF": responseHandler(replConfHandler),
	"SELECT":   replicationHandler(notImplementedHandler),
	"COMMAND":  replicationHandler(notImplementedHandler),
}

func replicationHandler(f handler) func(net.Conn, context.Context, []string) {
	return func(conn net.Conn, ctx context.Context, s []string) {
		f(ctx, s[1:])
	}
}

func handleReplicationConnection(conn net.Conn) {
	defer conn.Close()
	processConnection(conn, context.Background(), replicationHandlers)
}
