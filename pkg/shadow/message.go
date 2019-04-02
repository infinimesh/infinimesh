package shadow

import (
	"encoding/json"
	"time"
)

type FullDeviceStateMessage struct {
	Version   uint64
	State     json.RawMessage
	Timestamp time.Time
}

type DeltaDeviceStateMessage struct {
	Version   uint64
	Delta     json.RawMessage
	Timestamp time.Time
}
