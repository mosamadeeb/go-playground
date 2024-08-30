package main

func addEmailsToQueue(emails []string) chan string {
	emailsChan := make(chan string, len(emails))

	for _, v := range emails {
		emailsChan <- v
	}

	return emailsChan
}
