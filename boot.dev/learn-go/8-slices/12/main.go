package main

func indexOfFirstBadWord(msg []string, badWords []string) int {
	for i, word := range msg {
		for _, badWord := range badWords {
			if badWord == word {
				return i
			}
		}
	}

	return -1
}
