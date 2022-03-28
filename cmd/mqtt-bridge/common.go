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
	"io"
	"net"

	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"github.com/slntopp/mqtt-go/packet"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type VerifyDeviceFunc func(*devpb.Device) bool

func GetByFingerprintAndVerify(fingerprint []byte, cb VerifyDeviceFunc) (device *devpb.Device, err error) {
	resp, err := client.GetByFingerprint(internal_ctx, &devpb.GetByFingerprintRequest{
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
	log.Error("Closing connection on error", zap.Error(err))
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
	defer log.Info("Client disconnected", zap.String("client", connectPacket.ConnectPayload.ClientID))
	log.Info("Client connected", zap.String("client", connectPacket.ConnectPayload.ClientID))
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
		log.Info("Unsubbed from backchannel", zap.String("device", device.Uuid))
		ps.Unsub(backChannel)
	}()

	_, err := resp.WriteTo(c)
	if err != nil {
		log.Error("Failed to write Connection Acknowlegement", zap.Error(err))
		return
	}
	topicAliasPublishMap := make(map[string]int)

	token := device.GetToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer " + token)
	for {
		device, err = client.GetByToken(ctx, device)
		if err != nil {
			log.Error("Can't retrieve device status from registry", zap.Error(err))
		}
		if !device.Enabled {
			log.Debug("Device is disabled, disconnecting", zap.String("device", device.Uuid), zap.Bool("enabled", device.Enabled))
			_ = c.Close()
			break
		}
		p, err := packet.ReadPacket(c, connectPacket.VariableHeader.ProtocolLevel)

		if err != nil {
			if err == io.EOF {
				log.Info("Client closed connection", zap.String("client", connectPacket.ConnectPayload.ClientID))
			} else {
				log.Error("Failed to read packet", zap.Error(err))
			}
			_ = c.Close() // nolint: gosec
			break
		}

		switch p := p.(type) {
		case *packet.PingReqControlPacket:
			pong := packet.NewPingRespControlPacket()
			_, err := pong.WriteTo(c)
			if err != nil {
				log.Error("Failed to write Ping Response", zap.Error(err))
			}
		case *packet.PublishControlPacket:
			topicAliasPublishMap, err = handlePublish(p, c, device.Uuid, topicAliasPublishMap, int(connectPacket.VariableHeader.ProtocolLevel))
			if err != nil {
				log.Error("Failed to handle Publish", zap.Error(err))
			}

		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				log.Error("Failed to write Subscription Acknowlegement", zap.Error(err))
			}
			for _, sub := range p.Payload.Subscriptions {
				subTopic := TopicChecker(sub.Topic, device.Uuid)
				ps.AddSub(backChannel, subTopic)
				go handleBackChannel(c, device.Uuid, backChannel, connectPacket.VariableHeader.ProtocolLevel)
				log.Info("Added Subscription", zap.String("topic", subTopic), zap.String("device", device.Uuid))
			}
		case *packet.UnsubscribeControlPacket:
			response := packet.NewUnSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				log.Error("Failed to write Unsubscription Acknowlegement", zap.Error(err))
			}
			for _, unsub := range p.Payload.UnSubscriptions {
				ps.Unsub(backChannel, unsub.Topic)
				log.Info("Removed Subscription", zap.String("topic", unsub.Topic), zap.String("device", device.Uuid))
			}
		}
	}
}