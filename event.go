package main

type eventType int32

const (
	userJoined = iota
	userLeft
	message
)

type event struct {
	eventType
	content []byte
}

func newMessage(m []byte) event {
	return event{
		eventType: message,
		content:   m,
	}
}
