package shadow

import "encoding/json"

type DeviceState struct {
	Version int64
	State   json.RawMessage
}
