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
}

type handler func([]string) any

type transactionalHandler func(*Transaction, []string) any

func dummyTransactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) any {
		return f(s)
	}
}

func transactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) any {
		if *t != nil {
			t.addToQueue(func() any { return f(s) })
			return RespStr("QUEUED")
		} else {
			return f(s)
		}
	}
}

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
	var expiry *int = nil

	if intValue, err := strconv.Atoi(value.(string)); err == nil {
		value = intValue
	}

	if len(args) > 2 {
		option := args[2]
		optionValue := args[3]
		if strings.ToLower(option) == "px" {
			px, _ := strconv.Atoi(optionValue)
			expiry = &px
		}
	}

	db.set(key, value, expiry)
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
	db.set(key, newValue, nil)
	return newValue
}

func multiHandler(t *Transaction, args []string) any {
	*t = make(Transaction, 0)
	return RespStr("OK")
}

func execHandler(t *Transaction, args []string) any {
	if *t == nil {
		return errors.New("ERR EXEC without MULTI")
	}
	result := t.commit()
	*t = nil
	return result
}

func discardHandler(t *Transaction, args []string) any {
	if *t == nil {
		return errors.New("ERR DISCARD without MULTI")
	}
	*t = nil
	return RespStr("OK")
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
