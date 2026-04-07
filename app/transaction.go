package main

import "errors"

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

type transactionalHandler func(*Transaction, []string) any

func dummyTransactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) any {
		return f(s)
	}
}

func transactionHandler(f handler) transactionalHandler {
	return func(t *Transaction, s []string) any {
		if *t != nil {
			t.addToQueue(func() any { return f(s) })
			return RespStr("QUEUED")
		} else {
			return f(s)
		}
	}
}

func multiHandler(t *Transaction, args []string) any {
	*t = make(Transaction, 0)
	return RespStr("OK")
}

func execHandler(t *Transaction, args []string) any {
	if *t == nil {
		return errors.New("ERR EXEC without MULTI")
	}
	result := t.commit()
	*t = nil
	return result
}

func discardHandler(t *Transaction, args []string) any {
	if *t == nil {
		return errors.New("ERR DISCARD without MULTI")
	}
	*t = nil
	return RespStr("OK")
}
