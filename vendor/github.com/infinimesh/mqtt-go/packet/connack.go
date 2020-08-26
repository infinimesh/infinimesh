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
	"io"
)

type ConnAckControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader ConnAckVariableHeader
}

type ConnAckVariableHeader struct {
	SessionPresent bool
	ReasonCode     byte
}

func (p *ConnAckControlPacket) WriteTo(w io.Writer) (n int64, err error) {
	p.FixedHeader.RemainingLength = 2
	var nWritten int64
	nWritten, err = p.FixedHeader.WriteTo(w)
	n += nWritten
	if err != nil {
		return n, err
	}
	nWritten, err = p.VariableHeader.WriteTo(w)
	n += nWritten
	return n, err
}

func (c *ConnAckVariableHeader) WriteTo(w io.Writer) (n int64, err error) {
	buf := make([]byte, 2)

	buf[1] = c.ReasonCode

	bytesWritten, err := w.Write(buf)
	n += int64(bytesWritten)
	if err != nil {
		return
	}
	return
}
