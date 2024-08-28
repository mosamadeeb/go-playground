package main

import "strings"

// This is basically me solving a problem using another language's semantics
// Don't do this at home

// Rust unit struct
type unit struct{}

func countDistinctWords(messages []string) int {
	// We can implement a set by making a map without caring about the value type
	// We should either use bool (1 byte) or empty struct: struct{} (0 bytes)
	distinctWords := make(map[string]unit)

	for _, message := range messages {
		// Split by whitespace
		for _, word := range strings.Fields(message) {
			// Use lower case to avoid having multiple entries for differently capitalized words
			distinctWords[strings.ToLower(word)] = unit{}
		}
	}

	// Simple.
	return len(distinctWords)
}
