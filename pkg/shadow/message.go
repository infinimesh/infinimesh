package shadow

import "encoding/json"

type FullDeviceStateMessage struct {
	Version uint64
	State   json.RawMessage
}
