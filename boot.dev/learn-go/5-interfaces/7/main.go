package main

import "fmt"

func (e email) cost() int {
	if e.isSubscribed {
		return len(e.body) * 2
	} else {
		return len(e.body) * 5
	}
}

func (e email) format() string {
	subStatus := "Not Subscribed"

	if e.isSubscribed {
		subStatus = "Subscribed"
	}

	return fmt.Sprintf("'%s' | %s", e.body, subStatus)
}

type expense interface {
	cost() int
}

type formatter interface {
	format() string
}

type email struct {
	isSubscribed bool
	body         string
}
