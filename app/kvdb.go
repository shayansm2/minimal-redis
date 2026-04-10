package main

import (
	"sync"
	"time"
)

type KeyValueDB struct {
	kv       map[string]any
	mu       sync.RWMutex
	versions map[string]int
}

var db KeyValueDB

func init() {
	db = KeyValueDB{
		kv:       make(map[string]any),
		mu:       sync.RWMutex{},
		versions: make(map[string]int),
	}
}

func (db *KeyValueDB) set(key string, value any) {
	db.mu.Lock()
	defer db.mu.Unlock()
	version, found := db.versions[key]
	if !found {
		version = 0
	}
	db.kv[key] = value
	db.versions[key] = version + 1
}

func (db *KeyValueDB) expire(key string, ms int) {
	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		db.unset(key)
	}()
}

func (db *KeyValueDB) unset(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.kv, key)
	db.versions[key] = db.versions[key] + 1
}

func (db *KeyValueDB) get(key string) (any, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	value, found := (db.kv)[key]
	return value, found
}

func (db *KeyValueDB) keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, len(db.kv))
	i := 0
	for key := range db.kv {
		keys[i] = key
		i++
	}
	return keys
}

func (db *KeyValueDB) version(key string) int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.versions[key]
}
