package main

type messageToSend struct {
	message   string
	sender    user
	recipient user
}

type user struct {
	name   string
	number int
}

func (u user) isValid() bool {
	return u.name != "" && u.number != 0
}

func canSendMessage(mToSend messageToSend) bool {
	return mToSend.sender.isValid() && mToSend.recipient.isValid()
}
