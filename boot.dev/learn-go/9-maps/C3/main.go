package main

const (
	logDeleted  = "user deleted"
	logNotFound = "user not found"
	logAdmin    = "admin deleted"
)

type user struct {
	name   string
	number int
	admin  bool
}

func logAndDelete(users map[string]user, name string) (log string) {
	defer delete(users, name)

	user, ok := users[name]
	if !ok {
		return logNotFound
	}
	if user.admin {
		return logAdmin
	}

	return logDeleted
}
