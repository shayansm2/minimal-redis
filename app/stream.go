package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type ID struct {
	msTime int64
	seqNum int
}

func (id ID) Gt(other ID) bool {
	return (id.msTime > other.msTime) || (id.msTime == other.msTime && id.seqNum > other.seqNum)
}

func (id ID) Eq(other ID) bool {
	return id.msTime == other.msTime && id.seqNum == other.seqNum
}

func (id ID) toStr() string {
	return fmt.Sprintf("%d-%d", id.msTime, id.seqNum)
}

type Entry struct {
	id ID
	kv map[string]string
}

type Stream struct {
	ids []ID
	kvs []map[string]string
}

// var addToStreamEvents chan string

// func init() {
// 	addToStreamEvents = make(chan string)
// }

func getStream(key string) *Stream {
	s, found := db.get(key)
	if !found {
		return &Stream{}
	}
	return s.(*Stream)
}

func xAddHandler(ctx context.Context, args []string) any {
	key := args[0]
	stream := getStream(key)

	id, err := getNewValidId(stream, args[1])
	if err != nil {
		return err
	}

	kv := make(map[string]string)
	for i := 2; i < len(args); i += 2 {
		kv[args[i]] = args[i+1]
	}

	stream.ids = append(stream.ids, *id)
	stream.kvs = append(stream.kvs, kv)

	db.set(key, stream)
	// addToStreamEvents <- key

	return BulkStr(id.toStr())
}

// hard to refactor and change
func getNewValidId(stream *Stream, id string) (*ID, error) {
	lastId := ID{0, 0}
	if id == lastId.toStr() {
		return nil, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if len(stream.ids) > 0 {
		lastId = stream.ids[len(stream.ids)-1]
	}

	var msTime int64
	var seqNum int
	msTimeStr, seqNumStr, found := strings.Cut(id, "-")
	if !found {
		msTime = time.Now().UnixMilli()
		seqNumStr = "*"
	} else {
		ms, _ := strconv.Atoi(msTimeStr)
		msTime = int64(ms)
		seqNum, _ = strconv.Atoi(seqNumStr)
	}

	if msTime > lastId.msTime {
		if seqNumStr == "*" {
			seqNum = 0
		}
		return &ID{msTime, seqNum}, nil
	}

	if msTime < lastId.msTime {
		return nil, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if seqNumStr == "*" {
		seqNum = lastId.seqNum + 1
	}
	if seqNum > lastId.seqNum {
		return &ID{msTime, seqNum}, nil
	}
	return nil, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
}

func getQueryId(id string) ID {
	if !strings.Contains(id, "-") {
		msTime, _ := strconv.Atoi(id)
		return ID{int64(msTime), 0}
	}
	m, s, _ := strings.Cut(id, "-")
	msTime, _ := strconv.Atoi(m)
	seqNum, _ := strconv.Atoi(s)
	return ID{int64(msTime), seqNum}
}

func xRangeHandler(ctx context.Context, args []string) any {
	key := args[0]
	stream := getStream(key)
	var startIdx, endIdx int

	if args[1] == "-" {
		startIdx = 0
	} else {
		id := getQueryId(args[1])
		startIdx = utils.GreaterEqualIndex(stream.ids, id)
	}

	if args[2] == "+" {
		endIdx = len(stream.ids)
	} else {
		id := getQueryId(args[2])
		endIdx = min(utils.LessEqualIndex(stream.ids, id)+1, len(stream.ids))
	}

	return toStreamResponse(stream.ids[startIdx:endIdx], stream.kvs[startIdx:endIdx])
}

func xReadHandler(ctx context.Context, args []string) any {
	keyCount := (len(args) - 1) / 2
	result := make([]any, keyCount)
	for i := 1; i <= keyCount; i++ {
		key := args[i]
		stream := getStream(key)
		id := getQueryId(args[keyCount+i])
		idx := utils.GreaterThanIndex(stream.ids, id)
		result[i-1] = []any{BulkStr(key), toStreamResponse(stream.ids[idx:], stream.kvs[idx:])}
	}
	return result
}

func toStreamResponse(ids []ID, kvs []map[string]string) []any {
	result := make([]any, len(ids))
	for i, id := range ids {
		kv := kvs[i]
		kvArray := make([]BulkStr, len(kv)*2)
		j := 0
		for k, v := range kv {
			kvArray[j], kvArray[j+1] = BulkStr(k), BulkStr(v)
			j += 2
		}
		result[i] = []any{BulkStr(id.toStr()), kvArray}
	}
	return result
}
