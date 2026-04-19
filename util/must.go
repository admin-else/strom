package util

func MustOk[T any](v T, ok bool) T {
	if !ok {
		panic("not ok")
	}
	return v
}

func MustT[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
