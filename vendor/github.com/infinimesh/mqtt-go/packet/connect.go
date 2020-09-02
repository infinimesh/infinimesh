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
	"errors"
	"fmt"
	"io"
)

type ConnectProperties struct {
	PropertyLength         int //variable header properties length
	RecieveMaximumValue    int //limits the number of QoS 1 and QoS 2 Pub at Client - default 65,535
	MaximumPacketSize      int //represents max packet size client accepts
	SessionExpiryInterval  int //sesion expiry interval
	TopicAliasMaximumValue int //max num of topic alias accepted by client
	RequestResponseInfo    int //0 = no response info in CONNACK
	RequestProblemInfo     int //0 = no reason string in CONNACK
}

type ConnectFlags struct {
	UserName   bool
	Password   bool
	WillRetain bool
	WillQoS    int // 2 bytes actually
	WillFlag   bool
	CleanStart bool
}

type ConnectControlPacket struct {
	FixedHeader    FixedHeader
	VariableHeader ConnectVariableHeader
	ConnectPayload ConnectPayload
}

type ConnectVariableHeader struct {
	ProtocolName      string
	ProtocolLevel     byte
	ConnectFlags      ConnectFlags
	KeepAlive         int
	ConnectProperties ConnectProperties
}

type ConnectPayload struct {
	ClientID string
}

func getConnectVariableHeader(r io.Reader) (hdr ConnectVariableHeader, len int, err error) {
	// Protocol name
	protocolName, n, err := getProtocolName(r)
	len += n
	fmt.Printf("After protocolName Connect calculated length n %v\n", len)
	if err != nil {
		return hdr, 0, err
	}
	hdr.ProtocolName = protocolName

	if hdr.ProtocolName != "MQTT" && hdr.ProtocolName != "MQIsdp" {
		return hdr, 0, fmt.Errorf("Invalid protocol: %v", hdr.ProtocolName)
	}

	// Get Proto level
	protocolLevelBytes := make([]byte, 1)
	n, err = r.Read(protocolLevelBytes)
	len += n
	fmt.Printf("After protocolLevelBytes Connect calculated length n %v\n", len)
	if err != nil {
		return
	}
	hdr.ProtocolLevel = protocolLevelBytes[0]

	// Get Flags
	connectFlagsByte := make([]byte, 1)
	n, err = r.Read(connectFlagsByte)
	if n != 1 {
		return hdr, len, errors.New("Failed to read flags byte")
	}
	len += n
	if err != nil {
		return
	}

	hdr.ConnectFlags.UserName = connectFlagsByte[0]&128 == 1
	hdr.ConnectFlags.Password = connectFlagsByte[0]&64 == 1
	hdr.ConnectFlags.WillRetain = connectFlagsByte[0]&32 == 1
	hdr.ConnectFlags.WillFlag = connectFlagsByte[0]&4 == 1
	hdr.ConnectFlags.CleanStart = connectFlagsByte[0]&2 == 1

	keepAliveByte := make([]byte, 2)
	n, err = r.Read(keepAliveByte)
	len += n
	fmt.Printf("After keep alive Connect calculated length n %v\n", len)
	if err != nil {
		return hdr, len, errors.New("Could not read keepalive byte")
	}
	if n != 2 {
		return hdr, len, errors.New("Could not read enough keepalive bytes")
	}

	hdr.KeepAlive = int(binary.BigEndian.Uint16(keepAliveByte))

	// TODO Will QoS
	if connectFlagsByte[0]&16 == 1 && connectFlagsByte[0]&8 == 1 {
		hdr.ConnectFlags.WillQoS = 3
	} else if connectFlagsByte[0]&8 == 1 {
		hdr.ConnectFlags.WillQoS = 2
	} else {
		hdr.ConnectFlags.WillQoS = 1
	}

	//reading variable header properties length
	propertiesLength := make([]byte, 1)
	n, err = r.Read(propertiesLength)
	len += n
	if err != nil {
		return hdr, len, errors.New("Could not read properties length")
	}
	fmt.Printf("optional properties length %v and propertiesLength= %v\n ", n, propertiesLength)
	hdr.ConnectProperties.PropertyLength = int(propertiesLength[0])
	if hdr.ConnectProperties.PropertyLength < 1 {
		fmt.Printf("No optional properties added")
	} else {
		len += hdr.ConnectProperties.PropertyLength
		hdr, _ = readConnectProperties(r, hdr)
	}

	fmt.Printf("Connect calculated length n %v\n", len)
	return
}

func readConnectProperties(r io.Reader, hdr ConnectVariableHeader) (ConnectVariableHeader, error) {
	connectProperties := make([]byte, hdr.ConnectProperties.PropertyLength)
	propertiesLength, err := io.ReadFull(r, connectProperties)
	if err != nil {
		return hdr, err
	}
	if propertiesLength != hdr.ConnectProperties.PropertyLength {
		return hdr, errors.New("Connect Properties length incorrect")
	}
	for propertiesLength > 1 {
		connectPropertyID := int(connectProperties[0])
		if connectPropertyID == RECIEVE_MAXIMUM_ID {
			recieveMaximum := connectProperties[1 : RECIEVE_MAXIMUM_LENGTH+1]
			hdr.ConnectProperties.RecieveMaximumValue = int(binary.BigEndian.Uint16(recieveMaximum))
			connectProperties = connectProperties[RECIEVE_MAXIMUM_LENGTH+1 : propertiesLength]
			propertiesLength -= RECIEVE_MAXIMUM_LENGTH + 1
		} else if connectPropertyID == MAXIMUM_PACKET_SIZE_ID {
			maxPacketSize := connectProperties[1 : MAXIMUM_PACKET_SIZE_LENGTH+1]
			hdr.ConnectProperties.MaximumPacketSize = int(binary.BigEndian.Uint64(maxPacketSize))
			connectProperties = connectProperties[MAXIMUM_PACKET_SIZE_LENGTH+1 : propertiesLength]
			propertiesLength -= MAXIMUM_PACKET_SIZE_LENGTH + 1
		} else if connectPropertyID == SESSION_EXPIRY_INTERVAL_ID {
			SessionExpiryInterval := connectProperties[1 : SESSION_EXPIRY_INTERVAL_LENGTH+1]
			hdr.ConnectProperties.SessionExpiryInterval = int(binary.BigEndian.Uint64(SessionExpiryInterval))
			connectProperties = connectProperties[SESSION_EXPIRY_INTERVAL_LENGTH+1 : propertiesLength]
			propertiesLength -= SESSION_EXPIRY_INTERVAL_LENGTH + 1
		} else if connectPropertyID == TOPIC_ALIAS_MAXIMUM_ID {
			topicAliasMaximum := connectProperties[1 : TOPIC_ALIAS_MAXIMUM_LENGTH+1]
			hdr.ConnectProperties.TopicAliasMaximumValue = int(binary.BigEndian.Uint16(topicAliasMaximum))
			connectProperties = connectProperties[TOPIC_ALIAS_MAXIMUM_LENGTH+1 : propertiesLength]
			propertiesLength -= TOPIC_ALIAS_MAXIMUM_LENGTH + 1
		} else if connectPropertyID == REQUEST_RESPONSE_INFORMATION_ID {
			resquestResponseInfo := connectProperties[1]
			hdr.ConnectProperties.RequestResponseInfo = int(resquestResponseInfo)
			connectProperties = connectProperties[REQUEST_RESPONSE_INFORMATION_LENGTH:propertiesLength]
			propertiesLength -= REQUEST_RESPONSE_INFORMATION_LENGTH + 1
		} else if connectPropertyID == REQUEST_PROBLEM_INFORMATION_ID {
			resquestResponseInfo := connectProperties[1]
			hdr.ConnectProperties.RequestResponseInfo = int(resquestResponseInfo)
			connectProperties = connectProperties[REQUEST_PROBLEM_INFORMATION_LENGTH:propertiesLength]
			propertiesLength -= REQUEST_PROBLEM_INFORMATION_LENGTH + 1
		} else {
			fmt.Printf("%v Connect Property is not supported yet..", connectProperties[0])
			propertiesLength = 0
		}
	}
	return hdr, nil
}

func readConnectPayload(r io.Reader, len int) (ConnectPayload, error) {
	fmt.Printf("readConnectPayload len %v\n", len)
	payloadBytes := make([]byte, len)
	n, err := io.ReadFull(r, payloadBytes)
	fmt.Printf("payloadBytes len n %v and payload %v\n ", n, payloadBytes)
	// TODO set upper limit for payload
	// TODO only stream it
	if err != nil {
		return ConnectPayload{}, err
	}
	if n != len {
		return ConnectPayload{}, errors.New("Payload length incorrect")
	}

	// CONNECT MUST have the client id
	// REGEX 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
	// MAY allow more than that, but this must be possible

	// Client Identifier, Will Topic, Will Message, User Name, Password

	// TODO am besten so viel einlesen wie moeglich, und dann reslicen / reader zusammenstecken

	clientIDLengthBytes := payloadBytes[:2]
	fmt.Printf("clientIDLengthBytes = %v\n", clientIDLengthBytes)
	clientIDLength := int(clientIDLengthBytes[0])
	clientID := string(payloadBytes[1 : 1+clientIDLength])

	if clientIDLength == 0 {
		clientIDLength = int(binary.BigEndian.Uint16(clientIDLengthBytes))
		clientID = string(payloadBytes[2 : 2+clientIDLength])
	}
	fmt.Printf("clientIDLength = %v\n", clientIDLength)
	//clientID := string(payloadBytes[2 : 2+clientIDLength])
	return ConnectPayload{
		ClientID: clientID,
	}, nil

}
