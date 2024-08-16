package main

func (u User) SendMessage(message string, messageLength int) (string, bool) {
	if messageLength > u.MessageCharLimit {
		return "", false
	}

	return message, true
}

// don't touch below this line

type User struct {
	Name string
	Membership
}

type Membership struct {
	Type             membershipType
	MessageCharLimit int
}

type membershipType string

const (
	TypeStandard membershipType = "standard"
	TypePremium  membershipType = "premium"
)

func newUser(name string, membershipType membershipType) User {
	membership := Membership{Type: membershipType}
	if membershipType == TypeStandard {
		membership.MessageCharLimit = 100
	} else if membershipType == TypePremium {
		membership.MessageCharLimit = 1000
	}
	return User{Name: name, Membership: membership}
}
