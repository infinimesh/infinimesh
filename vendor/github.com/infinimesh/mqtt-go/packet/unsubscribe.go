package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type UnsubscribeUserProperty struct {
	key   string
	value string
}

type UnsubscribeProperties struct {
	PropertyLength int //1 byte
	UserProperty   UnsubscribeUserProperty
}

type UnsubscribeControlPacket struct {
	// Bits 3,2,1 and 0 of the fixed header of the SUBSCRIBE Control Packet are reserved and MUST be set to 0,0,1 and 0 respectively. The Server MUST treat any other value as malformed and close the Network Connection [MQTT-3.8.1-1].
	// TODO fail packet deserializing when this is not the case
	FixedHeader    FixedHeader
	VariableHeader UnsubscribeVariableHeader // 2 Bytes
	Payload        UnsubscribePayload
}

type UnsubscribeVariableHeader struct {
	PacketID              int // int16
	UnsubscribeProperties UnsubscribeProperties
}

type UnsubscribePayload struct {
	UnSubscriptions []Unsubscription
}

type Unsubscription struct {
	Topic string
	QoS   QosLevel
}

func readUnsubscribeVariableHeader(r io.Reader, protocolLevel byte) (n int, vh UnsubscribeVariableHeader, err error) {
	len := 0
	packetID, err := readUint16(r)
	len += 2
	if err != nil {
		return 0, UnsubscribeVariableHeader{}, err
	}
	vh.PacketID = packetID

	if int(protocolLevel) == 5 {
		propertyLength := make([]byte, 1)
		n, err = r.Read(propertyLength)
		len += n
		if err != nil {
			return
		}
		vh.UnsubscribeProperties.PropertyLength = int(propertyLength[0])
		if vh.UnsubscribeProperties.PropertyLength == 0 {
			fmt.Printf("No optional unsubscribe properties added")
		} else {
			len += vh.UnsubscribeProperties.PropertyLength
			vh, _ = readUnsubscribeProperties(r, vh)
		}
	}
	return len, vh, nil
}

func readUnsubscribeProperties(r io.Reader, vh UnsubscribeVariableHeader) (UnsubscribeVariableHeader, error) {
	unSubscribeProperties := make([]byte, vh.UnsubscribeProperties.PropertyLength)
	propertiesLength, err := io.ReadFull(r, unSubscribeProperties)
	if err != nil {
		return vh, err
	}
	if propertiesLength != vh.UnsubscribeProperties.PropertyLength {
		return vh, errors.New("Unsubscribe Properties length incorrect")
	}
	for propertiesLength > 1 {
		if unSubscribeProperties[0] == USER_PROPERTY_ID {
			unSubscribeProperties = unSubscribeProperties[1:]
			userPropertyKeyLength := int(binary.BigEndian.Uint16(unSubscribeProperties[0:2]))
			unSubscribeProperties = unSubscribeProperties[2:]

			vh.UnsubscribeProperties.UserProperty.key = string(unSubscribeProperties[0:userPropertyKeyLength])
			unSubscribeProperties = unSubscribeProperties[userPropertyKeyLength:]

			userPropertyValueLength := int(binary.BigEndian.Uint16(unSubscribeProperties[0:2]))
			unSubscribeProperties = unSubscribeProperties[2:]

			vh.UnsubscribeProperties.UserProperty.value = string(unSubscribeProperties[0:userPropertyValueLength])
			unSubscribeProperties = unSubscribeProperties[userPropertyValueLength:]
		} else {
			propertiesLength = 0
			fmt.Printf("No additional Unsubscribe properties added or supported")
		}
	}
	return vh, nil
}

func readUnsubscribePayload(r io.Reader, remainingLength int) (n int, payload UnsubscribePayload, err error) {
	for n < remainingLength {
		topicLength, err := readUint16(r)
		n += 2 // TODO get this info from readUint16, in case of errors it's maybe not exactly 2
		if err != nil {
			return n, UnsubscribePayload{}, err
		}

		topic := make([]byte, topicLength)
		bytesRead, err := io.ReadFull(r, topic)
		n += bytesRead
		if err != nil {
			return n, UnsubscribePayload{}, err
		}

		qos := make([]byte, 1)
		bytesRead, err = io.ReadFull(r, qos)
		n += bytesRead
		if err != nil {
			return n, UnsubscribePayload{}, err
		}

		unSub := Unsubscription{}
		unSub.Topic = string(topic)

		if qos[0]&252 > 0 {
			return n, UnsubscribePayload{}, errors.New("Invalid Unsubscribe payload. Reserved bits of QoS are non-zero")
		}

		if qos[0]&1 > 0 && qos[0]&2 > 0 {
			return n, UnsubscribePayload{}, errors.New("Invalid QoS level in payload. It is not allowed to set both bits")
		}

		if qos[0]&1 > 0 {
			unSub.QoS = QoSLevelAtLeastOnce
		} else if qos[0]&2 > 0 {
			unSub.QoS = QoSLevelExactlyOnce
		} else {
			unSub.QoS = QoSLevelNone
		}
		payload.UnSubscriptions = append(payload.UnSubscriptions, unSub)
	}
	return
}
