package main

import "strings"

type sms struct {
	id      string
	content string
	tags    []string
}

func tagMessages(messages []sms, tagger func(sms) []string) []sms {
	for i, m := range messages {
		messages[i].tags = tagger(m)
	}

	return messages
}

func tagger(msg sms) []string {
	// This is a good alternative to writing make([]string, 0)
	tags := []string{}

	if strings.Contains(msg.content, "urgent") || strings.Contains(msg.content, "Urgent") {
		tags = append(tags, "Urgent")
	}

	if strings.Contains(msg.content, "sale") || strings.Contains(msg.content, "Sale") {
		tags = append(tags, "Promo")
	}

	return tags
}
