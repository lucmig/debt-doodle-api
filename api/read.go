package api

// PongMessage - struct for pong messages
type PongMessage struct {
	Name  string
	Value string
}

// Pong - returns a pong
func Pong() PongMessage {
	return PongMessage{Name: "ping", Value: "pong"}
}
