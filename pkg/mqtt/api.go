//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package mqtt

type IncomingMessage struct {
	ProtoLevel   int
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

type Message struct {
	Topic string
	Data  map[string]interface{}
}

type Payload struct {
	Timestamp string
	Message   []Message
}
