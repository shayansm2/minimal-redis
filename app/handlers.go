package main

import (
	"errors"
	"strconv"
	"strings"
)

type handler func([]string) (string, error)

func errorHandler(f handler) func([]string) string {
	return func(s []string) string {
		respResult, err := f(s)
		if err != nil {
			return toRespError(err)
		}
		return respResult
	}
}

var handlers = map[string]func([]string) string{
	"ping": errorHandler(pingHandler),
	"echo": errorHandler(echoHandler),
	"set":  errorHandler(setHandler),
	"get":  errorHandler(getHandler),
	"incr": errorHandler(incrHandler),
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
