package must

// NotFail panics if the error is not nil, returns res otherwise.
//
// Use that function only for static initialization, test code, or code that "can't" fail.
// When in doubt, don't.
func NotFail[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

// NoError panics if the error is not nil.
//
// Use that function only for static initialization, test code, or code that "can't" fail.
// When in doubt, don't.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}
