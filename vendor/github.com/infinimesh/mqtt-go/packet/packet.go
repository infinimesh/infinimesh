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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ControlPacketType byte

type qosLevel int

// MQTT Quality of Service levels
const (
	QoSLevelNone         qosLevel = 0
	QoSLevelAtLeastOnce  qosLevel = 1
	QoSLevelExactyleOnce qosLevel = 2
)

// Control Packet types
const (
	CONNECT     = 1
	CONNACK     = 2
	PUBLISH     = 3
	PUBACK      = 4
	PUBREC      = 5
	PUBREL      = 6
	PUBCOMP     = 7
	SUBSCRIBE   = 8
	SUBACK      = 9
	UNSUBSCRIBE = 10
	UNSUBACK    = 11
	PINGREQ     = 12
	PINGRESP    = 13
	DISCONNECT  = 14
)

// FixedHeader is contained in every packet (thus, fixed). It consists of the
// Packet Type, Packet-specific Flags and the length of the rest of the message.
type FixedHeader struct {
	ControlPacketType ControlPacketType
	Flags             byte
	RemainingLength   int
}

type ControlPacket interface {
}

func getProtocolName(r io.Reader) (protocolName string, len int, err error) {
	protocolNameLengthBytes := make([]byte, 2)
	n, err := r.Read(protocolNameLengthBytes)
	len += n
	if err != nil {
		return "", len, errors.New("Failed to read length of protocolNameLengthBytes")
	}
	if n != 2 {

		return "", len, errors.New("Failed to read length of protocolNameLengthBytes, not enough bytes")
	}

	protocolNameLength := binary.BigEndian.Uint16(protocolNameLengthBytes)

	protocolNameBuffer := make([]byte, protocolNameLength)
	n, err = r.Read(protocolNameBuffer) // use ReadFull, its not guaranteed that we get enough out of a single read
	len += n
	if err != nil {
		return "", len, err
	}
	if n != int(protocolNameLength) {
		return "", len, err
	}

	return string(protocolNameBuffer), len, nil
}

func getFixedHeader(r io.Reader) (fh FixedHeader, err error) {
	buf := make([]byte, 1)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return FixedHeader{}, err
	}
	if n != 1 {
		return FixedHeader{}, errors.New("Failed to read MQTT Packet Control Type from Client Stream")
	}
	fh.ControlPacketType = ControlPacketType(buf[0] >> 4)
	fh.Flags = buf[0] & 15
	remainingLength, err := getRemainingLength(r) // Length VariableHeader + Payload
	if err != nil {
		return FixedHeader{}, err
	}
	fh.RemainingLength = remainingLength
	return
}

func ReadPacket(r io.Reader) (ControlPacket, error) {
	fh, err := getFixedHeader(r)
	if err != nil {
		return nil, err
	}

	// Ensure that we always read the remaining bytes
	bufRemaining := make([]byte, fh.RemainingLength)
	n, err := io.ReadFull(r, bufRemaining)
	if n != fh.RemainingLength {
		return nil, errors.New("short read")
	}
	if err != nil {
		return nil, err
	}

	remainingReader := bytes.NewBuffer(bufRemaining)

	return parseToConcretePacket(remainingReader, fh)
}

func parseToConcretePacket(remainingReader io.Reader, fh FixedHeader) (ControlPacket, error) {
	switch fh.ControlPacketType {
	case CONNECT:
		vh, variableHeaderSize, err := getConnectVariableHeader(remainingReader)
		if err != nil {
			return nil, err
		}
		payloadLength := fh.RemainingLength - variableHeaderSize

		cp, err := readConnectPayload(remainingReader, payloadLength)
		if err != nil {
			return nil, err
		}

		packet := &ConnectControlPacket{
			FixedHeader:    fh,
			VariableHeader: vh,
			ConnectPayload: cp,
		}

		return packet, nil
	case PUBLISH:
		flags, err := interpretPublishHeaderFlags(fh.Flags)
		if err != nil {
			return nil, err
		}

		vh, vhLength, err := readPublishVariableHeader(remainingReader, flags)
		if err != nil {
			return nil, err
		}

		payload, err := readPublishPayload(remainingReader, fh.RemainingLength-vhLength)
		if err != nil {
			return nil, err
		}

		packet := &PublishControlPacket{
			FixedHeader:      fh,
			FixedHeaderFlags: flags,
			VariableHeader:   vh,
			Payload:          payload,
		}
		return packet, nil

	case DISCONNECT:
		fmt.Println("Client disconnected")
		return nil, errors.New("Client disconnected")
	default:
		return nil, errors.New("Could not determine a specific control packet type")
	}

}

// starts with variable header

// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718023
func getRemainingLength(r io.Reader) (remaining int, err error) {
	// max 4 times / 4 rem. len.
	multiplier := 1
	for i := 0; i < 4; i++ {
		b := make([]byte, 1)
		n, err := r.Read(b)
		valueThisTime := int(b[0] & 127)
		remaining += valueThisTime * multiplier
		if err != nil {
			return remaining, err
		}
		if n != 1 {
			return 0, errors.New("Failed to get rem len")
		}

		multiplier *= 128
		moreBytes := b[0] & 128 // get only most significant bit
		if moreBytes == 0 {
			break
		}
	}
	return
}

func serializeRemainingLength(w io.Writer, len int) (n int, err error) {
	stuffToWrite := make([]byte, 0)
	for {
		encodedByte := byte(len % 128)
		len = len / 128

		if len > 0 {
			encodedByte |= 128 //set topmost bit to true because we
			//still have stuff to write
			stuffToWrite = append(stuffToWrite, encodedByte)
		} else {
			stuffToWrite = append(stuffToWrite, encodedByte)
			break
		}
	}
	return w.Write(stuffToWrite)
}

func (fh *FixedHeader) WriteTo(w io.Writer) (n int64, err error) {
	remainingLength := fh.RemainingLength
	b := byte(fh.ControlPacketType) << 4

	// Flags must be < 16
	b |= fh.Flags

	bytesWritten, err := w.Write([]byte{b})
	n += int64(bytesWritten)
	if err != nil {
		return
	}

	bytesWritten, err = serializeRemainingLength(w, remainingLength)
	n += int64(bytesWritten)
	return

}

func readUint16(r io.Reader) (len int, err error) {
	buf := make([]byte, 2)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return
	}
	if n != 2 {
		return n, errors.New("Couldnt read required 2 bytes for string length")
	}
	return int(binary.BigEndian.Uint16(buf)), nil
}
