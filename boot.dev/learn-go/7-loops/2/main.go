package main

func maxMessages(thresh int) int {
	cost := 0

	var numMessages int
	for numMessages = 0; cost < thresh; numMessages++ {
		cost += 100 + numMessages
	}

	if cost == thresh {
		// We managed to exactly fit the threshold
		return numMessages
	}

	// We can't send the last message
	return numMessages - 1
}
