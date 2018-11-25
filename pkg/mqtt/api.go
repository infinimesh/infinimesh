package mqtt

type MQTTBridgeData struct {
	SourceTopic  string
	SourceDevice string
	Data         []byte
}
