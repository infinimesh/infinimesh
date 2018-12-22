package mqtt

type IncomingMessage struct {
	SourceTopic  string
	SourceDevice string
	Data         []byte
}

// TBD: MQTT Subsystem will never be aware of the type/content of the message.
// Currently this is a deliberate design choice. Is it optimaL?
type OutgoingMessage struct {
	DeviceID string // "Target" device; does not necessarily have to be the connected device (e.g. sub-device)
	SubPath  string // Should not start with "/"
	Data     []byte
}
