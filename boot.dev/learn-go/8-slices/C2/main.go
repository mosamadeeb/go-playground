package main

import "unicode"

func isValidPassword(password string) bool {
	if len(password) < 5 || len(password) > 12 {
		return false
	}

	uppercaseFound := false
	digitFound := false

	for _, c := range password {
		if unicode.IsUpper(c) {
			uppercaseFound = true
		}

		if unicode.IsDigit(c) {
			digitFound = true
		}
	}

	return uppercaseFound && digitFound
}
