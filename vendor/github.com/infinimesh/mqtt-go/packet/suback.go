package packet

import (
	"bytes"
	"encoding/binary"
	"io"
)

type SubAckProperties struct {
	propertiesLength int
}

type SubAckControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader SubAckVariableHeader
	Payload        SubAckPayload
}

type SubAckVariableHeader struct {
	PacketID         uint16
	SubAckProperties SubAckProperties
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

func NewSubAck(packetID uint16, protocolLevel byte, returnCodes []byte) *SubAckControlPacket {
	if int(protocolLevel) == 5 {
		return &SubAckControlPacket{
			FixedHeader: FixedHeader{
				ControlPacketType: SUBACK,
				RemainingLength:   3 /* length of VH */ + len(returnCodes),
			},
			VariableHeader: SubAckVariableHeader{
				PacketID: packetID,
				SubAckProperties: SubAckProperties{
					propertiesLength: 1,
				},
			},
			Payload: SubAckPayload{
				ReturnCodes: returnCodes,
			},
		}
	}
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
	n, err = io.Copy(w, bytes.NewReader(b))
	if vh.SubAckProperties.propertiesLength != 0 {
		propertyLength := make([]byte, 1)
		nWritten, _ := w.Write(propertyLength)
		n += int64(nWritten)
	}
	return n, err
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
