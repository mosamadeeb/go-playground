package main

type membershipType string

const (
	TypeStandard membershipType = "standard"
	TypePremium  membershipType = "premium"
)

// don't touch above this line

type Membership struct {
	Type             membershipType
	MessageCharLimit int
}

type User struct {
	Membership
	Name string
}

func newUser(name string, membershipType membershipType) User {
	messageCharLimit := 100
	if membershipType == TypePremium {
		messageCharLimit = 1000
	}

	return User{
		Membership: Membership{
			Type:             membershipType,
			MessageCharLimit: messageCharLimit,
		},
		Name: name,
	}
}
