package main

import "time"

type DB map[string]string

var db DB

func init() {
	db = make(DB)
}

func (db *DB) set(key, value string, expiry *int) {
	(*db)[key] = value
	if expiry != nil {
		go func() {
			time.Sleep(time.Duration(*expiry) * time.Millisecond)
			db.unset(key)
		}()
	}
}

func (db *DB) unset(key string) {
	delete(*db, key)
}

func (db *DB) get(key string) (string, bool) {
	value, found := (*db)[key]
	return value, found
}
