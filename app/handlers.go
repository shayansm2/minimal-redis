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
	value := args[1]
	var expiry *int = nil

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
	return bulkStringEncode(value)
}

func incrHandler(args []string) string {
	key := args[0]
	value, _ := db.get(key)
	numericValue, _ := strconv.Atoi(value)
	newNumericValue := numericValue + 1
	newValue := strconv.Itoa(newNumericValue)
	db.set(key, newValue, nil)
	return respIntegerEncode(newNumericValue)
}
