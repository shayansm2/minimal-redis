package main

import (
	"context"
	"errors"
)

type Watcher map[string]int

func NewWatcher() Watcher {
	return make(Watcher)
}

func (l *Watcher) watch(name string) {
	(*l)[name] = db.version(name)
}

func (l *Watcher) unwatch(name string) {
	delete(*l, name)
}

func (l *Watcher) reset() {
	*l = NewWatcher()
}

func (l *Watcher) hasChanged() bool {
	for key, version := range *l {
		if version != db.version(key) {
			return true
		}
	}
	return false
}

func watchHandler(ctx context.Context, args []string) any {
	if len(args) < 1 {
		return errors.New("ERR not enough args provided")
	}
	t := ctx.Value(TransactionContextKey).(*Transaction)
	if t.isInTransaction() {
		return errors.New("ERR WATCH inside MULTI is not allowed")
	}
	watcher := ctx.Value(WatcherContextKey).(*Watcher)
	for _, key := range args {
		watcher.watch(key)
	}

	return RespStr("OK")
}

func unwatchHandler(ctx context.Context, args []string) any {
	watcher := ctx.Value(WatcherContextKey).(*Watcher)
	watcher.reset()
	return RespStr("OK")
}
