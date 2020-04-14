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
