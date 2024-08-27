package main

func getMessageWithRetries(primary, secondary, tertiary string) ([3]string, [3]int) {
	messages := [3]string{primary, secondary, tertiary}
	var costs [3]int

	costs[0] = len(messages[0])

	for i := 1; i < 3; i++ {
		costs[i] = len(messages[i]) + costs[i-1]
	}

	return messages, costs
}
