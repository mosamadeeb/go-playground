package main

func countConnections(groupSize int) int {
	numConnections := 0

	// A group with 1 person has no connections
	for i := 2; i <= groupSize; i++ {
		numConnections += i - 1
	}

	return numConnections
}
