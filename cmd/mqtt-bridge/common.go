/*
Copyright © 2021-2023 Infinite Devices GmbH

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

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/infinimesh/infinimesh/pkg/mqtt/metrics"
	"github.com/infinimesh/infinimesh/pkg/pubsub"
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

// LogErrorAndClose - Logs Error, Sends Acknowlegement(ACK) packet and Closes the connection
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

// HandleConn - note: Connection is expected to be valid & legitimate at this point
func HandleConn(c net.Conn, connectPacket *packet.ConnectControlPacket, device *devpb.Device) {
	log := log.Named(device.GetUuid()).Named(connectPacket.ConnectPayload.ClientID)

	defer log.Debug("Client disconnected")

	log.Debug(
		"Client connected", zap.String("device", device.Uuid),
		zap.Int("protocol_level", int(connectPacket.VariableHeader.ProtocolLevel)),
		zap.String("protocol", connectPacket.VariableHeader.ProtocolName),
		zap.Int("QoS", connectPacket.VariableHeader.ConnectFlags.WillQoS),
	)
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

	metrics.ActiveConnectionsTotal.Inc()
	ps.TryPub(&pb.Shadow{
		Device: device.Uuid,
		Connection: &pb.ConnectionState{
			Connected: true,
			Timestamp: timestamppb.Now(),
			ClientId:  connectPacket.ConnectPayload.ClientID,
		},
	}, "mqtt.incoming")

	defer func() {
		metrics.ActiveConnectionsTotal.Dec()
		ps.TryPub(&pb.Shadow{
			Device: device.Uuid,
			Connection: &pb.ConnectionState{
				Connected: false,
				Timestamp: timestamppb.Now(),
				ClientId:  connectPacket.ConnectPayload.ClientID,
			},
		}, "mqtt.incoming")
	}()

	token := device.GetToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	for {
		device, err = client.GetByToken(ctx, device)
		if err != nil {
			log.Warn("Can't retrieve device status from registry", zap.Error(err))
		}
		if !device.Enabled {
			log.Debug("Device is disabled, disconnecting", zap.String("device", device.Uuid), zap.Bool("enabled", device.Enabled))
			err = c.Close()
			log.Warn("Error closing connection", zap.Error(err))
			break
		}

		p, err := packet.ReadPacket(c, connectPacket.VariableHeader.ProtocolLevel)
		if err != nil {
			if err == io.EOF {
				log.Debug("Client closed connection", zap.String("client", connectPacket.ConnectPayload.ClientID))
			} else {
				log.Warn("Failed to read packet", zap.Error(err))
			}
			if err := c.Close(); err != nil {
				log.Warn("Couldn't close connection", zap.Error(err))
			}
			return
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
				Connection: &pb.ConnectionState{
					Connected: true,
					Timestamp: timestamppb.Now(),
				},
			}
			ps.TryPub(payload, "mqtt.incoming")

			// _, err := packet.NewPubAckControlPacket(uint16(p.VariableHeader.PacketID)).WriteTo(c)
			// if err != nil {
			// 	log.Warn("Failed to write Publish Acknowlegement", zap.Error(err))
			// }

		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				log.Warn("Failed to write Subscription Acknowlegement", zap.Error(err))
			}

			for _, sub := range p.Payload.Subscriptions {
				ps.AddSub(backChannel, "mqtt.outgoing/"+device.Uuid)
				go handleBackChannel(log, backChannel, c, sub.Topic, connectPacket.VariableHeader.ProtocolLevel, func() {
					ps.TryPub(&pb.Shadow{
						Device: device.Uuid,
						Connection: &pb.ConnectionState{
							Connected: true,
							Timestamp: timestamppb.Now(),
						},
					}, "mqtt.incoming")
				})
				log.Debug("Added Subscription", zap.String("topic", sub.Topic), zap.String("device", device.Uuid))
			}

			go func() {
				if shadow != nil {
					r, err := shadow.Get(ctx, &pb.GetRequest{Pool: []string{device.Uuid}})
					if err != nil || len(r.GetShadows()) == 0 {
						return
					}
					state := r.GetShadows()[0]
					if state.Desired != nil {
						ps.TryPub(state, "mqtt.outgoing/"+device.Uuid)
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
				log.Debug("Removed Subscription", zap.String("topic", unsub.Topic), zap.String("device", device.Uuid))
			}
		}
	}
}

func handleBackChannel(log *zap.Logger, ch chan interface{}, c net.Conn, topic string, protocolLevel byte, connected func()) {
	defer log.Debug("BackChannel handler closed")
	var ts int64 = 0
	for msg := range ch {
		shadow := msg.(*pb.Shadow)
		log.Debug("Received message", zap.String("topic", topic), zap.String("device", shadow.Device))
		if shadow.Desired == nil || shadow.Desired.Timestamp == nil {
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
			log.Error("Failed to write packet", zap.Error(err))
			return
		}

		connected()
	}
}

func unsub[T chan any](ps pubsub.PubSub, ch chan any) {
	go ps.Unsub(ch)

	for range ch {
	}
}
