package main

type DB map[string]string

var db DB

func init() {
	db = make(DB)
}

func (db *DB) set(key, value string) {
	(*db)[key] = value
}

func (db *DB) get(key string) (string, bool) {
	value, found := (*db)[key]
	return value, found
}
