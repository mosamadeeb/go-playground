package main

func sum(nums ...int) int {
	currentSum := 0

	for i := 0; i < len(nums); i++ {
		currentSum += nums[i]
	}

	return currentSum
}
