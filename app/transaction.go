package main

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
