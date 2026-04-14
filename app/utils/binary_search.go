package utils

type Comparable[T any] interface {
	Gt(other T) bool
	Eq(other T) bool
}

func EqualIndex[T Comparable[T]](array []T, find T) int {
	s, e := 0, len(array)-1
	for s <= e {
		m := (s + e) / 2
		i := (array)[m]
		if i.Eq(find) {
			return m
		}
		if i.Gt(find) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return s - 1
}

func LessEqualIndex[T Comparable[T]](array []T, find T) int {
	s, e := 0, len(array)-1
	for s <= e {
		m := (s + e) / 2
		i := (array)[m]
		if i.Eq(find) {
			return m
		}
		if i.Gt(find) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return s - 1
}

func GreaterEqualIndex[T Comparable[T]](array []T, find T) int {
	s, e := 0, len(array)-1
	for s <= e {
		m := (s + e) / 2
		i := (array)[m]
		if i.Eq(find) {
			return m
		}
		if i.Gt(find) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return e + 1
}

func GreaterThanIndex[T Comparable[T]](array []T, find T) int {
	s, e := 0, len(array)-1
	for s <= e {
		m := (s + e) / 2
		i := (array)[m]
		if i.Eq(find) {
			return m + 1
		}
		if i.Gt(find) {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	return e + 1
}
