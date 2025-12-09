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

func IsType[T any](a any) bool {
	switch a.(type) {
	case T:
		return true
	default:
		return false
	}
}
