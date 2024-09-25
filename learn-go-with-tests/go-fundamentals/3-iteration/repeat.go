package iteration

func Repeat(character string, repeatCount int) string {
	var repeated string
	for range repeatCount {
		repeated += character
	}

	return repeated
}
