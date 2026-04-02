package main

import "time"

type KeyValueDB map[string]string

var db KeyValueDB

func init() {
	db = make(KeyValueDB)
}

func (db *KeyValueDB) set(key, value string, expiry *int) {
	(*db)[key] = value
	if expiry != nil {
		go func() {
			time.Sleep(time.Duration(*expiry) * time.Millisecond)
			db.unset(key)
		}()
	}
}

func (db *KeyValueDB) unset(key string) {
	delete(*db, key)
}

func (db *KeyValueDB) get(key string) (string, bool) {
	value, found := (*db)[key]
	return value, found
}
