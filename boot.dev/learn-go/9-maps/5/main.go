package main

func getCounts(messagedUsers []string, validUsers map[string]int) {
	for _, name := range messagedUsers {
		if _, ok := validUsers[name]; ok {
			validUsers[name]++
		}
	}
}
