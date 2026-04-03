package main

import (
	"strconv"
	"strings"
)

type handler func([]string) string

var handlers = map[string]handler{
	"ping": pingHandler,
	"echo": echoHandler,
	"set":  setHandler,
	"get":  getHandler,
	"incr": incrHandler,
}

func pingHandler(args []string) string {
	return respSimpleStringEncode("PONG")
}

func echoHandler(args []string) string {
	str := args[0]
	return bulkStringEncode(str)
}

func setHandler(args []string) string {
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
	return respSimpleStringEncode("OK")
}

func getHandler(args []string) string {
	key := args[0]
	value, found := db.get(key)
	if !found {
		return NullBulkString
	}

	if intVal, ok := value.(int); ok {
		// are you sure? can get respond with resp int?
		return respIntegerEncode(intVal)
	}
	return bulkStringEncode(value.(string))
}

func incrHandler(args []string) string {
	key := args[0]
	value, _ := db.get(key)
	newValue := value.(int) + 1
	db.set(key, newValue, nil)
	return respIntegerEncode(newValue)
}
