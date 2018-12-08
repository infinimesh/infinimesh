package shadow

import "encoding/json"

type FullDeviceStateMessage struct {
	Version int64
	State   json.RawMessage
}
