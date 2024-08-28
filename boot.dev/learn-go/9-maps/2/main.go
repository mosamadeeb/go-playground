package main

import "errors"

func deleteIfNecessary(users map[string]user, name string) (deleted bool, err error) {
	u, ok := users[name]

	if !ok {
		return false, errors.New("not found")
	}

	if !u.scheduledForDeletion {
		return false, nil
	}

	delete(users, name)
	return true, nil
}

type user struct {
	name                 string
	number               int
	scheduledForDeletion bool
}
