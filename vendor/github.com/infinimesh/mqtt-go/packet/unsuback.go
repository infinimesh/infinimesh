package packet

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UnSubAckProperties struct {
	propertiesLength int
}

type UnSubAckControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader UnSubAckVariableHeader
	Payload        UnSubAckPayload
}

type UnSubAckVariableHeader struct {
	PacketID           uint16
	UnSubAckProperties UnSubAckProperties
}

type UnSubAckPayload struct {
	ReturnCodes []byte
}

func NewUnSubAck(packetID uint16, protocolLevel byte, returnCodes []byte) *UnSubAckControlPacket {
	if int(protocolLevel) == 5 {
		return &UnSubAckControlPacket{
			FixedHeader: FixedHeader{
				ControlPacketType: UNSUBACK,
				RemainingLength:   3 /* length of VH */ + len(returnCodes),
			},
			VariableHeader: UnSubAckVariableHeader{
				PacketID: packetID,
				UnSubAckProperties: UnSubAckProperties{
					propertiesLength: 1,
				},
			},
			Payload: UnSubAckPayload{
				ReturnCodes: returnCodes,
			},
		}
	}
	return &UnSubAckControlPacket{
		FixedHeader: FixedHeader{
			ControlPacketType: UNSUBACK,
			RemainingLength:   2 /* length of VH */ + len(returnCodes),
		},
		VariableHeader: UnSubAckVariableHeader{
			PacketID: packetID,
		},
		Payload: UnSubAckPayload{
			ReturnCodes: returnCodes,
		},
	}
}

// TODO deserializing

func (vh *UnSubAckVariableHeader) WriteTo(w io.Writer) (n int64, err error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, vh.PacketID)
	n, err = io.Copy(w, bytes.NewReader(b))
	if vh.UnSubAckProperties.propertiesLength > 0 {
		propertyLength := make([]byte, 1)
		nWritten, _ := w.Write(propertyLength)
		n += int64(nWritten)
	}
	return n, err
}

func (p *UnSubAckControlPacket) WriteTo(w io.Writer) (n int64, err error) {
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
