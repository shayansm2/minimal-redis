package main

import (
	"time"
)

type KeyValueDB map[string]any

var db KeyValueDB

func init() {
	db = make(KeyValueDB)
}

func (db *KeyValueDB) set(key string, value any) {
	(*db)[key] = value
}

func (db *KeyValueDB) expire(key string, ms int) {
	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		db.unset(key)
	}()
}

func (db *KeyValueDB) unset(key string) {
	delete(*db, key)
}

func (db *KeyValueDB) get(key string) (any, bool) {
	value, found := (*db)[key]
	return value, found
}

func (db *KeyValueDB) keys() []string {
	keys := make([]string, len(*db))
	i := 0
	for key := range *db {
		keys[i] = key
		i++
	}
	return keys
}
