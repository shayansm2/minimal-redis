package main

import (
	"context"
	"errors"
)

type Transaction []func() any

func NewTransaction() Transaction {
	return make(Transaction, 0)
}

func (t *Transaction) addToQueue(f func() any) {
	*t = append(*t, f)
}

func (t *Transaction) isInTransaction() bool {
	return *t != nil
}

func (t *Transaction) commit() []string {
	result := make([]string, len(*t))
	for i, f := range *t {
		res := f()
		result[i], _ = encode(res)
	}
	*t = nil
	return result
}

func (t *Transaction) abort() {
	*t = nil
}

func transactionHandler(f handler) handler {
	return func(ctx context.Context, s []string) any {
		t := ctx.Value(TransactionContextKey).(*Transaction)
		if t.isInTransaction() {
			t.addToQueue(func() any { return f(ctx, s) })
			return RespStr("QUEUED")
		} else {
			return f(ctx, s)
		}
	}
}

func multiHandler(ctx context.Context, args []string) any {
	t := ctx.Value(TransactionContextKey).(*Transaction)
	*t = NewTransaction()
	return RespStr("OK")
}

func execHandler(ctx context.Context, args []string) any {
	t := ctx.Value(TransactionContextKey).(*Transaction)
	if !t.isInTransaction() {
		return errors.New("ERR EXEC without MULTI")
	}
	watcher := ctx.Value(WatcherContextKey).(*Watcher)
	defer func() { watcher.reset() }()
	if watcher.hasChanged() {
		t.abort()
		var nullArr []string
		return nullArr
	}
	result := t.commit()
	return result
}

func discardHandler(ctx context.Context, args []string) any {
	t := ctx.Value(TransactionContextKey).(*Transaction)
	if !t.isInTransaction() {
		return errors.New("ERR DISCARD without MULTI")
	}
	t.abort()
	return RespStr("OK")
}
