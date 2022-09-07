/*
Copyright © 2021-2022 Infinite Devices GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/cskr/pubsub"
	structpb "github.com/golang/protobuf/ptypes/struct"
	devpb "github.com/infinimesh/proto/node/devices"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/slntopp/mqtt-go/packet"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	log.Warn("Closing connection on error", zap.Error(err))
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
	log.Info("Client connected", zap.String("device", device.Uuid), zap.String("client", connectPacket.ConnectPayload.ClientID))
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
	defer unsub(ps, backChannel)

	_, err := resp.WriteTo(c)
	if err != nil {
		log.Warn("Failed to write Connection Acknowlegement", zap.Error(err))
		return
	}

	token := device.GetToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	for {
		device, err = client.GetByToken(ctx, device)
		if err != nil {
			log.Warn("Can't retrieve device status from registry", zap.Error(err))
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
				log.Warn("Failed to read packet", zap.Error(err))
			}
			_ = c.Close() // nolint: gosec
			break
		}

		switch p := p.(type) {
		case *packet.PingReqControlPacket:
			pong := packet.NewPingRespControlPacket()
			_, err := pong.WriteTo(c)
			if err != nil {
				log.Warn("Failed to write Ping Response", zap.Error(err))
			}
		case *packet.PublishControlPacket:
			var data structpb.Struct
			err = data.UnmarshalJSON(p.Payload)
			if err != nil {
				log.Warn("Failed to handle Publish", zap.Error(err))
				continue
			}
			payload := &pb.Shadow{
				Device: device.Uuid,
				Reported: &pb.State{
					Timestamp: timestamppb.Now(),
					Data:      &data,
				},
			}
			ps.Pub(payload, "mqtt.incoming")
		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				log.Warn("Failed to write Subscription Acknowlegement", zap.Error(err))
			}

			for _, sub := range p.Payload.Subscriptions {
				ps.AddSub(backChannel, "mqtt.outgoing/"+device.Uuid)
				go handleBackChannel(backChannel, c, sub.Topic, connectPacket.VariableHeader.ProtocolLevel)
				log.Info("Added Subscription", zap.String("topic", sub.Topic), zap.String("device", device.Uuid))
			}

			go func() {
				if shadow != nil {
					r, err := shadow.Get(ctx, &pb.GetRequest{Pool: []string{device.Uuid}})
					if err != nil || len(r.GetShadows()) == 0 {
						return
					}
					state := r.GetShadows()[0]
					if state.Desired != nil {
						ps.Pub(state, "mqtt.outgoing/"+device.Uuid)
					}
				}
			}()
		case *packet.UnsubscribeControlPacket:
			response := packet.NewUnSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				log.Warn("Failed to write Unsubscription Acknowlegement", zap.Error(err))
			}
			for _, unsub := range p.Payload.UnSubscriptions {
				ps.Unsub(backChannel, "mqtt.outgoing/"+device.Uuid)
				log.Info("Removed Subscription", zap.String("topic", unsub.Topic), zap.String("device", device.Uuid))
			}
		}
	}
}

func handleBackChannel(ch chan interface{}, c net.Conn, topic string, protocolLevel byte) {
	var ts int64 = 0
	for msg := range ch {
		shadow := msg.(*pb.Shadow)
		log.Debug("Received message", zap.String("topic", topic), zap.String("device", shadow.Device))
		if shadow.Desired == nil {
			log.Debug("Skipping empty Desired state")
			continue
		}
		if shadow.Desired.Timestamp.Seconds < ts {
			log.Debug("Skipping message", zap.String("topic", topic), zap.String("device", shadow.Device))
			continue
		}
		payload, err := shadow.Desired.Data.MarshalJSON()
		if err != nil {
			log.Warn("Failed to marshal shadow", zap.Error(err))
			continue
		}
		p := packet.NewPublish(topic, 0, payload, protocolLevel)
		_, err = p.WriteTo(c)
		if err != nil {
			panic(err)
		}
	}
}

func unsub[T chan any](ps *pubsub.PubSub, ch chan any) {
	go ps.Unsub(ch)

	for range ch {
	}
}
