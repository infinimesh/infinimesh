package shadow

import (
	"encoding/json"
	"time"
)

type DeviceStateMessage struct {
	Version   uint64
	State     json.RawMessage
	Timestamp time.Time
}
