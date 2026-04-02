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
}

func pingHandler(args []string) string {
	return respStringEncode("PONG")
}

func echoHandler(args []string) string {
	str := args[0]
	return bulkEncode(str)
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
	return respStringEncode("OK")
}

func getHandler(args []string) string {
	key := args[0]
	value, found := db.get(key)
	if !found {
		return NullBulkString
	}
	return bulkEncode(value)
}
