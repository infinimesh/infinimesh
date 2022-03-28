//--------------------------------------------------------------------------
// Copyright 2018-2022 infinimesh
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

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"github.com/slntopp/mqtt-go/packet"
	"google.golang.org/grpc/metadata"
)

type VerifyDeviceFunc func(*devpb.Device) bool

func GetByFingerprintAndVerify(fingerprint []byte, cb VerifyDeviceFunc) (device *devpb.Device, err error) {
	resp, err := client.GetByFingerprint(context.Background(), &devpb.GetByFingerprintRequest{
		Fingerprint: fingerprint,
	})
	if err != nil {
		return nil, err
	}

	if cb(resp) {
		return resp, nil
	}

	return nil, errors.New("not found")
}

// Log Error, Send Acknowlegement(ACK) packet and Close the connection
// ACK Packet needs to be sent to prevent MQTT Client sending CONN packets further
func LogErrorAndClose(c net.Conn, err error) {
	fmt.Printf("Closing connection on error: %v\n", err)
	resp := packet.ConnAckControlPacket{
		FixedHeader: packet.FixedHeader{
			ControlPacketType: packet.CONNACK,
		},
		VariableHeader: packet.ConnAckVariableHeader{},
	}
	resp.WriteTo(c)
	c.Close()
}

// Connection is expected to be valid & legitimate at this point
func HandleConn(c net.Conn, connectPacket *packet.ConnectControlPacket, device *devpb.Device) {
	defer fmt.Println("Client disconnected ", connectPacket.ConnectPayload.ClientID)

	fmt.Printf("Client with ID %v connected!\n", connectPacket.ConnectPayload.ClientID)
	// TODO ignore/compare this ID with the given ID from the verify function

	//TODO : MQTT CONNACK Properties need to add here
	resp := packet.ConnAckControlPacket{
		FixedHeader: packet.FixedHeader{
			ControlPacketType: packet.CONNACK,
		},
		VariableHeader: packet.ConnAckVariableHeader{},
	}

	if len(connectPacket.ConnectPayload.ClientID) <= 0 {
		resp.VariableHeader.ConnAckProperties.AssignedClientID = device.Uuid
	}
	// Only open Back-channel after conn packet was received

	// Create empty subscription
	backChannel := ps.Sub()
	go handleBackChannel(c, device.Uuid, backChannel, connectPacket.VariableHeader.ProtocolLevel)
	defer func() {
		fmt.Printf("Unsubbed channel %v\n", device.Uuid)
		ps.Unsub(backChannel)
	}()

	_, err := resp.WriteTo(c)
	if err != nil {
		fmt.Println("Failed to write ConnAck. Closing connection.")
		return
	}
	topicAliasPublishMap := make(map[string]int)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer" + device.GetToken())
	for {

		// TODO: check device status(Enabeld/Disabled)

		device, err = client.GetByToken(ctx, device)
		if err != nil {
			fmt.Printf("device status doesn't exist in redis %v\n", err)
		} else {
			if !device.Enabled {
				_ = c.Close()
				break
			}
		}
		p, err := packet.ReadPacket(c, connectPacket.VariableHeader.ProtocolLevel)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client closed connection.\n")
			} else {
				fmt.Printf("Error while reading packet in client loop: %v. Disconnecting client.\n", err)
			}
			_ = c.Close() // nolint: gosec
			break
		}

		switch p := p.(type) {
		case *packet.PingReqControlPacket:
			pong := packet.NewPingRespControlPacket()
			_, err := pong.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write PingResp", err)
			}
		case *packet.PublishControlPacket:
			topicAliasPublishMap, err = handlePublish(p, c, device.Uuid, topicAliasPublishMap, int(connectPacket.VariableHeader.ProtocolLevel))
			if err != nil {
				fmt.Printf("Failed to handle Publish packet: %v.", err)
			}

		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write SubAck:", err)
			}
			for _, sub := range p.Payload.Subscriptions {
				subTopic := TopicChecker(sub.Topic, device.Uuid)
				ps.AddSub(backChannel, subTopic)
				go handleBackChannel(c, device.Uuid, backChannel, connectPacket.VariableHeader.ProtocolLevel)
				fmt.Println("Added Subscription", subTopic, device.Uuid)
			}
		case *packet.UnsubscribeControlPacket:
			response := packet.NewUnSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write UnSubAck:", err)
			}
			for _, unsub := range p.Payload.UnSubscriptions {
				ps.Unsub(backChannel, unsub.Topic)
				fmt.Println("Removed Subscription", unsub.Topic, device.Uuid)
			}
		}
	}
}