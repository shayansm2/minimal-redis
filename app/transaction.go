package main

import (
	"context"
	"errors"
)

type Transaction []func() any

func (t *Transaction) addToQueue(f func() any) {
	*t = append(*t, f)
}

func (t *Transaction) commit() []string {
	result := make([]string, len(*t))
	for i, f := range *t {
		res := f()
		result[i], _ = encode(res)
	}
	return result
}

func transactionHandler(f handler) handler {
	return func(ctx context.Context, s []string) any {
		t := ctx.Value("transaction").(*Transaction)
		if *t != nil {
			t.addToQueue(func() any { return f(ctx, s) })
			return RespStr("QUEUED")
		} else {
			return f(ctx, s)
		}
	}
}

func multiHandler(ctx context.Context, args []string) any {
	t := ctx.Value("transaction").(*Transaction)
	*t = make(Transaction, 0)
	return RespStr("OK")
}

func execHandler(ctx context.Context, args []string) any {
	t := ctx.Value("transaction").(*Transaction)
	if *t == nil {
		return errors.New("ERR EXEC without MULTI")
	}
	result := t.commit()
	*t = nil
	return result
}

func discardHandler(ctx context.Context, args []string) any {
	t := ctx.Value("transaction").(*Transaction)
	if *t == nil {
		return errors.New("ERR DISCARD without MULTI")
	}
	*t = nil
	return RespStr("OK")
}
