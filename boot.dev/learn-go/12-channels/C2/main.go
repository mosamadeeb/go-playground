package main

import "time"

// This article was very useful for this entire chapter
// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
func processMessages(messages []string) []string {
	// We set the length from the start instead of setting capacity and appending
	// This allows us to assign the processed messages according to index, so they will stay ordered
	// This actually isn't necessary for the test case
	processed := make([]string, len(messages))

	ch := make(chan struct{})

	for i, v := range messages {
		go func() {
			processed[i] = process(v)
			ch <- struct{}{}
		}()
	}

	// Ensure channel all messages are processed before returning
	for range len(messages) {
		<-ch
	}

	return processed
}

// don't touch below this line

func process(message string) string {
	time.Sleep(1 * time.Second)
	return message + "-processed"
}
