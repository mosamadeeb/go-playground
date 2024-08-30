package main

func getLast[T any](s []T) T {
	if len(s) == 0 {
		var defaultVal T
		return defaultVal
	}

	return s[len(s)-1]
}
