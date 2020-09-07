package packet

import "io"

type PingRespControlPacket struct {
	FixedHeader FixedHeader
}

func (p *PingRespControlPacket) WriteTo(w io.Writer) (n int64, err error) {
	return p.FixedHeader.WriteTo(w)
}

func NewPingRespControlPacket() *PingRespControlPacket {
	return &PingRespControlPacket{
		FixedHeader: FixedHeader{
			ControlPacketType: PINGRESP,
		},
	}
}
