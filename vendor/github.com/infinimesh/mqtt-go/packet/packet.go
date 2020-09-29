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

type QosLevel int

// MQTT Quality of Service levels
const (
	QoSLevelNone        QosLevel = 0
	QoSLevelAtLeastOnce QosLevel = 1
	QoSLevelExactlyOnce QosLevel = 2
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
const (
	SESSION_EXPIRY_INTERVAL_ID          = 17
	SESSION_EXPIRY_INTERVAL_LENGTH      = 4
	RECIEVE_MAXIMUM_ID                  = 33
	RECIEVE_MAXIMUM_LENGTH              = 2
	MAXIMUM_PACKET_SIZE_ID              = 39
	MAXIMUM_PACKET_SIZE_LENGTH          = 4
	TOPIC_ALIAS_MAXIMUM_ID              = 34
	TOPIC_ALIAS_MAXIMUM_LENGTH          = 2
	REQUEST_RESPONSE_INFORMATION_ID     = 25
	REQUEST_RESPONSE_INFORMATION_LENGTH = 1
	REQUEST_PROBLEM_INFORMATION_ID      = 23
	REQUEST_PROBLEM_INFORMATION_LENGTH  = 1
	TOPIC_ALIAS_ID                      = 35
	TOPIC_ALIAS_LENGTH                  = 2
	MESSAGE_EXPIRY_INTERVAL_ID          = 2
	MESSAGE_EXPIRY_INTERVAL_LENGTH      = 4
	RESPONSE_TOPIC_ID                   = 8
	RESPONSE_TOPIC_LENGTH               = 1
	CORRELATION_DATA_ID                 = 9
	CORRELATION_DATA_LENGTH             = 1
	USER_PROPERTY_ID                    = 38
	USER_PROPERTY_LENGTH                = 1
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

func getProtocolLevel(r io.Reader) (protocolLevel byte, len int, err error) {
	// Get Proto level
	protocolLevelBytes := make([]byte, 1)
	len, err = r.Read(protocolLevelBytes)
	if err != nil {
		return protocolLevelBytes[0], len, err
	}
	return protocolLevelBytes[0], len, nil
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

func ReadPacket(r io.Reader, protocolLevel byte) (ControlPacket, error) {
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

	return parseToConcretePacket(remainingReader, fh, protocolLevel)
}

// nolint: gocyclo
func parseToConcretePacket(remainingReader io.Reader, fh FixedHeader, protocolLevel byte) (ControlPacket, error) {
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

		vh, vhLength, err := readPublishVariableHeader(remainingReader, flags, protocolLevel)
		if err != nil {
			return nil, err
		}

		payload, err := readPublishPayload(remainingReader, fh.RemainingLength-vhLength)
		fmt.Printf("Publish payload :%v\n", payload)
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
	case SUBSCRIBE:
		vhLen, vh, err := readSubscribeVariableHeader(remainingReader, protocolLevel)
		if err != nil {
			return nil, err
		}

		_, payload, err := readSubscribePayload(remainingReader, fh.RemainingLength-vhLen)
		if err != nil {
			return nil, err
		}

		packet := &SubscribeControlPacket{
			FixedHeader:    fh,
			VariableHeader: vh,
			Payload:        payload,
		}
		return packet, nil
	case PINGREQ:
		return &PingReqControlPacket{FixedHeader: fh}, nil
	case UNSUBSCRIBE:
		vhLen, vh, err := readUnsubscribeVariableHeader(remainingReader, protocolLevel)
		if err != nil {
			return nil, err
		}

		_, payload, err := readUnsubscribePayload(remainingReader, fh.RemainingLength-vhLen)
		if err != nil {
			return nil, err
		}

		packet := &UnsubscribeControlPacket{
			FixedHeader:    fh,
			VariableHeader: vh,
			Payload:        payload,
		}
		fmt.Printf("Client has unsubscribed %v", fh.ControlPacketType)
		return packet, nil
	case DISCONNECT:
		fmt.Println("Client disconnected")
		return nil, errors.New("Client disconnected")
	default:
		return nil, fmt.Errorf("Unknown control packet type: %v", fh.ControlPacketType)
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

// Allocating here everytime is super inefficient, better pass a byte
// slice
// TODO return number of bytes read
func readUint16(r io.Reader) (result int, err error) {
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
