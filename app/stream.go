package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ID struct {
	msTime int64
	seqNum int
}

func (id ID) gt(other ID) bool {
	return (id.msTime > other.msTime) || (id.msTime == other.msTime && id.seqNum > other.seqNum)
}

func (id ID) lt(other ID) bool {
	return (id.msTime < other.msTime) || (id.msTime == other.msTime && id.seqNum < other.seqNum)
}

func (id ID) eq(other ID) bool {
	return id.msTime == other.msTime && id.seqNum == other.seqNum
}

func (id ID) toStr() string {
	return fmt.Sprintf("%d-%d", id.msTime, id.seqNum)
}

type Entry struct {
	id ID
	kv map[string]string
}

type Stream []Entry

func leIdx(stream *Stream, id ID) int {
	s, e := 0, len(*stream)-1
	for s <= e {
		m := (s + e) / 2
		i := (*stream)[m].id
		if i.eq(id) {
			return m
		}
		if i.gt(id) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return s - 1
}

func geIdx(stream *Stream, id ID) int {
	s, e := 0, len(*stream)-1
	for s <= e {
		m := (s + e) / 2
		i := (*stream)[m].id
		if i.eq(id) {
			return m
		}
		if i.gt(id) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return e + 1
}

func gdIdx(stream *Stream, id ID) int {
	s, e := 0, len(*stream)-1
	for s <= e {
		m := (s + e) / 2
		i := (*stream)[m].id
		if i.eq(id) {
			return m + 1
		}
		if i.gt(id) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return e + 1
}

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

	id, err := getValidNewId(stream, args[1])
	if err != nil {
		return err
	}

	kv := make(map[string]string)
	for i := 2; i < len(args); i += 2 {
		kv[args[i]] = args[i+1]
	}

	*stream = append(*stream, Entry{id: *id, kv: kv})

	db.set(key, stream)

	return BulkStr(id.toStr())
}

// hard to refactor and change
func getValidNewId(stream *Stream, id string) (*ID, error) {
	lastId := ID{0, 0}
	if id == lastId.toStr() {
		return nil, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if len(*stream) > 0 {
		lastId = (*stream)[len(*stream)-1].id
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
		startIdx = geIdx(stream, getQueryId(args[1]))
	}

	if args[2] == "+" {
		endIdx = len(*stream)
	} else {
		endIdx = min(leIdx(stream, getQueryId(args[2]))+1, len(*stream))
	}

	return toStreamResponse((*stream)[startIdx:endIdx])
}

func xReadHandler(ctx context.Context, args []string) any {
	keyCount := (len(args) - 1) / 2
	result := make([]any, keyCount)
	for i := 1; i <= keyCount; i++ {
		key := args[i]
		stream := getStream(key)
		id := getQueryId(args[keyCount+i])
		idx := gdIdx(stream, id)
		result[i-1] = []any{BulkStr(key), toStreamResponse((*stream)[idx:])}
	}
	return result
}

func toStreamResponse(stream []Entry) []any {
	result := make([]any, len(stream))
	for i, entry := range stream {
		kvArray := make([]BulkStr, len(entry.kv)*2)
		j := 0
		for k, v := range entry.kv {
			kvArray[j], kvArray[j+1] = BulkStr(k), BulkStr(v)
			j += 2
		}
		result[i] = []any{BulkStr(entry.id.toStr()), kvArray}
	}
	return result
}
