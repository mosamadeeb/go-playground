package main

func concurrentFib(n int) []int {
	fibVals := make([]int, 0, n)

	ch := make(chan int)
	go fibonacci(n, ch)

	// This will block on each iteration until there is a value to receive
	// It will also exit the loop once the channel is closed
	for x := range ch {
		fibVals = append(fibVals, x)
	}

	return fibVals
}

// don't touch below this line

func fibonacci(n int, ch chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		ch <- x
		x, y = y, x+y
	}
	close(ch)
}
