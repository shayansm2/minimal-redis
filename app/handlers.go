package main

func pingHandler() string {
	return respStringEncode("PONG")
}

func echoHandler(str string) string {
	return bulkEncode(str)
}

func setHandler(key, value string) string {
	db.set(key, value)
	return respStringEncode("OK")
}

func getHandler(key string) string {
	value, found := db.get(key)
	if !found {
		return NullBulkString
	}
	return bulkEncode(value)
}
