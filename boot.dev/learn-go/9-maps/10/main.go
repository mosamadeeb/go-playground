package main

func getNameCounts(names []string) map[rune]map[string]int {
	runeToName := make(map[rune]map[string]int)

	for _, name := range names {
		firstRune := []rune(name)[0]

		nameMap, ok := runeToName[firstRune]
		if !ok {
			nameMap = make(map[string]int)
			runeToName[firstRune] = nameMap
		}

		if _, ok := nameMap[name]; !ok {
			nameMap[name] = 0
		}

		nameMap[name]++
	}

	return runeToName
}
