package main

import (
	"errors"
	"strconv"
	"strings"
)

var handlers = map[string]func(*Transaction, []string) string{
	"PING":    responseHandler(dummyTransactionHandler(pingHandler)),
	"ECHO":    responseHandler(dummyTransactionHandler(echoHandler)),
	"SET":     responseHandler(transactionHandler(setHandler)),
	"GET":     responseHandler(transactionHandler(getHandler)),
	"INCR":    responseHandler(transactionHandler(incrHandler)),
	"MULTI":   responseHandler(multiHandler),
	"EXEC":    responseHandler(execHandler),
	"DISCARD": responseHandler(discardHandler),
	"CONFIG":  responseHandler(dummyTransactionHandler(configHandler)),
	"KEYS":    responseHandler(dummyTransactionHandler(keysHandler)),
	"RPUSH":   responseHandler(transactionHandler(rPushHandler)),
	"LPUSH":   responseHandler(transactionHandler(lPushHandler)),
	"LRANGE":  responseHandler(transactionHandler(lRangeHandler)),
	"LLEN":    responseHandler(transactionHandler(lLenHandler)),
	"LPOP":    responseHandler(transactionHandler(lPopHandler)),
}

type handler func([]string) any

func responseHandler(f transactionalHandler) func(*Transaction, []string) string {
	return func(t *Transaction, s []string) string {
		result := f(t, s)
		encoded, err := encode(result)
		if err != nil {
			return toRespError(err)
		}
		return encoded
	}
}

func pingHandler(args []string) any {
	return RespStr("PONG")
}

func echoHandler(args []string) any {
	if len(args) == 0 {
		return errors.New("not enough args provided")
	}
	return BulkStr(args[0])
}

func setHandler(args []string) any {
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

func getHandler(args []string) any {
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

func incrHandler(args []string) any {
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

func configHandler(args []string) any {
	if len(args) < 2 {
		return errors.New("ERR not enough args provided")
	}
	key := args[1]
	value := configs[key]
	return []BulkStr{BulkStr(key), BulkStr(value)}
}

func keysHandler(args []string) any {
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
