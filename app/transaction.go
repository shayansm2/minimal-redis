package main

type Transaction []func() (any, error)

func (t *Transaction) addToQueue(f func() (any, error)) {
	*t = append(*t, f)
}

// todo handling errors
func (t *Transaction) commit() []string {
	result := make([]string, len(*t))
	for i, f := range *t {
		res, _ := f()
		result[i], _ = encode(res)
	}
	return result
}
