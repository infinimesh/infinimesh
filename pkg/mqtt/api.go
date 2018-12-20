package mqtt

type IncomingMessage struct {
	SourceTopic  string
	SourceDevice string
	Data         []byte
}

type OutgoingMessage struct {
	Topic string
	Data  []byte
	// Device is not necessary; anyone who is allowed to access given topic may get it
}
