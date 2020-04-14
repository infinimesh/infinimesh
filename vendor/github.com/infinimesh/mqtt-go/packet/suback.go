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

package packet

import (
	"bytes"
	"encoding/binary"
	"io"
)

type SubAckControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader SubAckVariableHeader
	Payload        SubAckPayload
}

type SubAckVariableHeader struct {
	PacketID uint16
}

type SubAckPayload struct {
	ReturnCodes []byte
}

// Allowed return codes:

// 0x00 - Success - Maximum QoS 0
// 0x01 - Success - Maximum QoS 1
// 0x02 - Success - Maximum QoS 2
// 0x80 - Failure
const (
	ReturncodeSuccessQoS0 byte = 0x00
	ReturncodeSuccessQoS1 byte = 0x01
	ReturncodeSuccessQoS2 byte = 0x02
	ReturncodeFailure     byte = 0x80
)

func NewSubAck(packetID uint16, returnCodes []byte) *SubAckControlPacket {
	return &SubAckControlPacket{
		FixedHeader: FixedHeader{
			ControlPacketType: SUBACK,
			RemainingLength:   2 /* length of VH */ + len(returnCodes),
		},
		VariableHeader: SubAckVariableHeader{
			PacketID: packetID,
		},
		Payload: SubAckPayload{
			ReturnCodes: returnCodes,
		},
	}
}

// TODO deserializing

func (vh *SubAckVariableHeader) WriteTo(w io.Writer) (n int64, err error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, vh.PacketID)

	return io.Copy(w, bytes.NewReader(b))
}

func (p *SubAckControlPacket) WriteTo(w io.Writer) (n int64, err error) {
	written, err := p.FixedHeader.WriteTo(w)
	n += written
	if err != nil {
		return
	}

	written, err = p.VariableHeader.WriteTo(w)
	n += written
	if err != nil {
		return
	}

	wr, err := w.Write(p.Payload.ReturnCodes)
	n += int64(wr)
	if err != nil {
		return n, err
	}
	return
}
