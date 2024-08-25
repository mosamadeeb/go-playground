package main

func bulkSend(numMessages int) float64 {
	cost := 0.0

	for i := 0; i < numMessages; i++ {
		cost += 1.0 + float64(i)/100.0
	}

	return cost

	// Or, if we want to look smart
	// n * (n + 1) / 2
	// return float64(numMessages) + float64((numMessages-1)*(numMessages)/2)/100.0
}
