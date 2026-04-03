package main

import (
	"errors"
	"strconv"
	"strings"
)

var handlers = map[string]func(Transaction, []string) string{
	"ping":  errorHandler(dummyTransactionHandler(pingHandler)),
	"echo":  errorHandler(dummyTransactionHandler(echoHandler)),
	"set":   errorHandler(transactionHandler(setHandler)),
	"get":   errorHandler(transactionHandler(getHandler)),
	"incr":  errorHandler(transactionHandler(incrHandler)),
	"multi": errorHandler(multiHandler),
	"exec":  errorHandler(execHandler),
}

type handler func([]string) (string, error)

type transactionalHandler func(Transaction, []string) (string, error)

func dummyTransactionHandler(f handler) transactionalHandler {
	return func(t Transaction, s []string) (string, error) {
		return f(s)
	}
}

func transactionHandler(f handler) transactionalHandler {
	return func(t Transaction, s []string) (string, error) {
		if *t {
			return toRespSimpleString("QUEUED"), nil
		} else {
			return f(s)
		}
	}
}

func errorHandler(f transactionalHandler) func(Transaction, []string) string {
	return func(t Transaction, s []string) string {
		respResult, err := f(t, s)
		if err != nil {
			return toRespError(err)
		}
		return respResult
	}
}

func pingHandler(args []string) (string, error) {
	return toRespSimpleString("PONG"), nil
}

func echoHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("not enough args provided")
	}
	str := args[0]
	return toBulkString(str), nil
}

func setHandler(args []string) (string, error) {
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
	return toRespSimpleString("OK"), nil
}

func getHandler(args []string) (string, error) {
	key := args[0]
	value, found := db.get(key)
	if !found {
		return NullBulkString, nil
	}

	if _, ok := value.(int); ok {
		return toBulkString(strconv.Itoa(value.(int))), nil
	}
	return toBulkString(value.(string)), nil
}

func incrHandler(args []string) (string, error) {
	key := args[0]
	value, found := db.get(key)
	if !found {
		value = 0
	}
	if _, isString := value.(string); isString {
		return "", errors.New("ERR value is not an integer or out of range")
	}
	newValue := value.(int) + 1
	db.set(key, newValue, nil)
	return toRespInteger(newValue), nil
}

func multiHandler(t Transaction, args []string) (string, error) {
	*t = true
	return toRespSimpleString("OK"), nil
}

func execHandler(t Transaction, args []string) (string, error) {
	if !*t {
		return "", errors.New("ERR EXEC without MULTI")
	}
	*t = false
	return toRespArray([]string{}), nil
}
