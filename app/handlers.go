package main

import (
	"errors"
	"strconv"
	"strings"
)

var handlers = map[string]func(*Transaction, []string) string{
	"ping":  responseHandler(dummyTransactionHandler(pingHandler)),
	"echo":  responseHandler(dummyTransactionHandler(echoHandler)),
	"set":   responseHandler(transactionHandler(setHandler)),
	"get":   responseHandler(transactionHandler(getHandler)),
	"incr":  responseHandler(transactionHandler(incrHandler)),
	"multi": responseHandler(multiHandler),
	"exec":  responseHandler(execHandler),
}

type handler func([]string) (any, error)

type transactionalHandler func(*Transaction, []string) (any, error)

func dummyTransactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) (any, error) {
		return f(s)
	}
}

func transactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) (any, error) {
		if *t != nil {
			t.addToQueue(func() (any, error) { return f(s) })
			return RespStr("QUEUED"), nil
		} else {
			return f(s)
		}
	}
}

func responseHandler(f transactionalHandler) func(*Transaction, []string) string {
	return func(t *Transaction, s []string) string {
		result, err := f(t, s)
		if err != nil {
			return toRespError(err)
		}
		encoded, err := encode(result)
		if err != nil {
			return toRespError(err)
		}
		return encoded
	}
}

func pingHandler(args []string) (any, error) {
	return RespStr("PONG"), nil
}

func echoHandler(args []string) (any, error) {
	if len(args) == 0 {
		return "", errors.New("not enough args provided")
	}
	return BulkStr(args[0]), nil
}

func setHandler(args []string) (any, error) {
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
	return RespStr("OK"), nil
}

func getHandler(args []string) (any, error) {
	key := args[0]
	value, found := db.get(key)
	if !found {
		return nil, nil
	}

	if _, ok := value.(int); ok {
		return BulkStr(strconv.Itoa(value.(int))), nil
	}
	return BulkStr(value.(string)), nil
}

func incrHandler(args []string) (any, error) {
	key := args[0]
	value, found := db.get(key)
	if !found {
		value = 0
	}
	if _, isString := value.(string); isString {
		return 0, errors.New("ERR value is not an integer or out of range")
	}
	newValue := value.(int) + 1
	db.set(key, newValue, nil)
	return newValue, nil
}

func multiHandler(t *Transaction, args []string) (any, error) {
	*t = make(Transaction, 0)
	return RespStr("OK"), nil
}

func execHandler(t *Transaction, args []string) (any, error) {
	if *t == nil {
		return nil, errors.New("ERR EXEC without MULTI")
	}
	result := t.commit()
	*t = nil
	return result, nil
}
