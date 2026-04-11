package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	id string
	kv map[string]string
}

type Stream []Entry

func xAddHandler(ctx context.Context, args []string) any {
	key := args[0]
	id := args[1]

	kv := make(map[string]string)
	for i := 2; i < len(args); i += 2 {
		kv[args[i]] = args[i+1]
	}

	s, found := db.get(key)
	var stream *Stream
	if !found {
		stream = &Stream{}
	} else {
		stream = s.(*Stream)
	}

	id, err := getValidId(stream, id)
	if err != nil {
		return err
	}

	*stream = append(*stream, Entry{id: id, kv: kv})
	db.set(key, stream)
	return BulkStr(id)
}

func getValidId(stream *Stream, id string) (string, error) {
	lastId := "0-0"
	if id == lastId {
		return "", errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if len(*stream) > 0 {
		lastId = (*stream)[len(*stream)-1].id
	}
	msTime, seqNum, found := strings.Cut(id, "-")
	lastMsTime, lastSeqNum, _ := strings.Cut(lastId, "-")
	if !found {
		msTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
		seqNum = "*"
	}
	if msTime > lastMsTime {
		if seqNum == "*" {
			seqNum = "0"
		}
		return fmt.Sprintf("%s-%s", msTime, seqNum), nil
	}
	if msTime < lastMsTime {
		return "", errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}
	if seqNum == "*" {
		n, _ := strconv.Atoi(lastSeqNum)
		seqNum = strconv.Itoa(n + 1)
	}
	if seqNum > lastSeqNum {
		return fmt.Sprintf("%s-%s", msTime, seqNum), nil
	}
	return "", errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
}
