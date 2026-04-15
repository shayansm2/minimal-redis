package main

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

var handlers = map[string]func(context.Context, []string) string{
	"PING":    responseHandler(pingHandler),
	"ECHO":    responseHandler(echoHandler),
	"SET":     responseHandler(transactionHandler(setHandler)),
	"GET":     responseHandler(transactionHandler(getHandler)),
	"INCR":    responseHandler(transactionHandler(incrHandler)),
	"MULTI":   responseHandler(multiHandler),
	"EXEC":    responseHandler(execHandler),
	"DISCARD": responseHandler(discardHandler),
	"CONFIG":  responseHandler(configHandler),
	"KEYS":    responseHandler(keysHandler),
	"RPUSH":   responseHandler(transactionHandler(rPushHandler)),
	"LPUSH":   responseHandler(transactionHandler(lPushHandler)),
	"LRANGE":  responseHandler(transactionHandler(lRangeHandler)),
	"LLEN":    responseHandler(transactionHandler(lLenHandler)),
	"LPOP":    responseHandler(transactionHandler(lPopHandler)),
	"BLPOP":   responseHandler(transactionHandler(bLPopHandler)),
	"WATCH":   responseHandler(watchHandler),
	"UNWATCH": responseHandler(unwatchHandler),
	"TYPE":    responseHandler(typeHandler),
	"XADD":    responseHandler(xAddHandler),
	"XRANGE":  responseHandler(xRangeHandler),
	"XREAD":   responseHandler(xReadHandler),
	"INFO":    responseHandler(infoHandler),
}

type handler func(context.Context, []string) any

func responseHandler(f handler) func(context.Context, []string) string {
	return func(ctx context.Context, s []string) string {
		result := f(ctx, s)
		encoded, err := encode(result)
		if err != nil {
			return toRespError(err)
		}
		return encoded
	}
}

func pingHandler(ctx context.Context, args []string) any {
	return RespStr("PONG")
}

func echoHandler(ctx context.Context, args []string) any {
	if len(args) == 0 {
		return errors.New("not enough args provided")
	}
	return BulkStr(args[0])
}

func setHandler(ctx context.Context, args []string) any {
	key := args[0]
	var value any = args[1]

	if intValue, err := strconv.Atoi(value.(string)); err == nil {
		value = intValue
	}
	db.set(key, value)
	if len(args) > 2 {
		option := args[2]
		optionValue := args[3]
		if strings.ToLower(option) == "px" {
			px, _ := strconv.Atoi(optionValue)
			db.expire(key, px)
		}
	}

	return RespStr("OK")
}

func getHandler(ctx context.Context, args []string) any {
	key := args[0]
	value, found := db.get(key)
	if !found {
		return nil
	}

	if _, ok := value.(int); ok {
		return BulkStr(strconv.Itoa(value.(int)))
	}
	return BulkStr(value.(string))
}

func incrHandler(ctx context.Context, args []string) any {
	key := args[0]
	value, found := db.get(key)
	if !found {
		value = 0
	}
	if _, isString := value.(string); isString {
		return errors.New("ERR value is not an integer or out of range")
	}
	newValue := value.(int) + 1
	db.set(key, newValue)
	return newValue
}

func configHandler(ctx context.Context, args []string) any {
	if len(args) < 2 {
		return errors.New("ERR not enough args provided")
	}
	key := args[1]
	value := configs[key]
	return []BulkStr{BulkStr(key), BulkStr(value)}
}

func keysHandler(ctx context.Context, args []string) any {
	if len(args) < 1 {
		return errors.New("ERR not enough args provided")
	}
	keys := db.keys()
	bulkKeys := make([]BulkStr, len(keys))
	for i, key := range keys {
		bulkKeys[i] = BulkStr(key)
	}
	return bulkKeys
}

const TypeString = "string"
const TypeStream = "stream"
const TypeNone = "none"

func typeHandler(ctx context.Context, args []string) any {
	key := args[0]
	val, found := db.get(key)
	if !found {
		return RespStr(TypeNone)
	}
	switch val.(type) {
	case *Stream:
		return RespStr(TypeStream)
	default:
		return RespStr(TypeString)
	}
}
