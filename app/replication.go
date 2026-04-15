package main

import (
	"context"
	"fmt"
	"strings"
)

var role string

func init() {
	if getConfigs().get("replicaof", "") != "" {
		role = "slave"
	} else {
		role = "master"
	}
}

func infoHandler(ctx context.Context, args []string) any {
	info := []string{fmt.Sprintf(
		"role:%s", role),
		"master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		"master_repl_offset:0",
	}
	return BulkStr(strings.Join(info, "\n"))
}
