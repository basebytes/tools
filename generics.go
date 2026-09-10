package tools

func ToAny[T any](a any) T {
	var t T
	switch a.(type) {
	case T:
		return a.(T)
	default:
		return t
	}
}

func ToAnyE[T any](a any) (T, bool) {
	var t T
	switch a.(type) {
	case T:
		return a.(T), true
	default:
		return t, false
	}
}

func ToAnySlice[T, S any](a []S) []T {
	results := make([]T, len(a))
	for i, v := range a {
		results[i] = ToAny[T](v)
	}
	return results
}

func ToAnySliceE[T, S any](a []S) (results []T, ok bool) {
	results = make([]T, len(a))
	for i, v := range a {
		if results[i], ok = ToAnyE[T](v); !ok {
			break
		}
	}
	return
}

func IsType[T any](a any) bool {
	switch a.(type) {
	case T:
		return true
	default:
		return false
	}
}

func SymDiff[T comparable](a, b []T) (da, db []T) {
	if la := len(a); la == 0 {
		db = b
	} else if lb := len(b); lb == 0 {
		da = a
	} else if (la > lb && la/2 > lb) || (lb > la && lb/2 < la) {
		db, da = SymDiff[T](b, a)
	} else {
		m := make(map[T]int, la)
		db = make([]T, 0, lb)
		da = make([]T, 0, la)
		for _, v := range a {
			m[v] = 1
		}
		for _, v := range b {
			if _, ok := m[v]; ok {
				m[v]--
			} else {
				db = append(db, v)
			}
		}
		for _, v := range a {
			if m[v] > 0 {
				da = append(da, v)
			}
		}
	}
	return
}
