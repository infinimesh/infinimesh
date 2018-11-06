//--------------------------------------------------------------------------
// Copyright 2018 infinimesh, INC
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

package packet

import (
	"encoding/binary"
	"io"
)

type PubackControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader PubAckVariableHeader
}

type PubAckVariableHeader struct {
	PacketID uint16
}

func (vh *PubAckVariableHeader) WriteTo(w io.Writer) (n int64, err error) {
	packetID := make([]byte, 2)
	binary.BigEndian.PutUint16(packetID, vh.PacketID)

	bytesWritten, err := w.Write(packetID)
	n += int64(bytesWritten)
	if err != nil {
		return
	}
	return
}

func (p *PubackControlPacket) WriteTo(w io.Writer) (n int64, err error) {
	return p.VariableHeader.WriteTo(w)
}

func NewPubAckControlPacket(packetID uint16) *PubackControlPacket {
	return &PubackControlPacket{
		FixedHeader: FixedHeader{
			ControlPacketType: PUBACK,
			RemainingLength:   2,
		},
		VariableHeader: PubAckVariableHeader{
			PacketID: packetID,
		},
	}
}
