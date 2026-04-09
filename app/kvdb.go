package main

import (
	"sync"
	"time"
)

type KeyValueDB struct {
	kv map[string]any
	mu sync.RWMutex
}

var db KeyValueDB

func init() {
	db = KeyValueDB{
		kv: make(map[string]any),
		mu: sync.RWMutex{},
	}
}

func (db *KeyValueDB) set(key string, value any) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.kv[key] = value
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
