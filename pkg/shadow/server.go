//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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

package shadow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

//Server is a data strcuture for shadow server
type Server struct {
	shadowpb.UnimplementedShadowsServer

	Repo         Repo
	Producer     sarama.SyncProducer // Sync producer, we want to guarantee execution
	ProduceTopic string
	Log          *zap.Logger

	PubSub *pubsub.PubSub
}

//Get is a method to get a device state
func (s *Server) Get(context context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {

	log := s.Log.Named("Get State Controller")
	log.Debug("Function Invoked", zap.String("Device", request.Id))

	response = &shadowpb.GetResponse{
		Shadow: &shadowpb.Shadow{},
	}

	// TODO fetch device from registry, 404 if not found

	reportedState, err := s.Repo.GetReported(request.Id)
	if err != nil {
		reportedState.ID = request.Id
		reportedState.State = DeviceStateMessage{
			Version: uint64(0),
			State:   json.RawMessage([]byte("{}")),
		}
	}

	desiredState, err := s.Repo.GetDesired(request.Id)
	if err != nil {
		desiredState.ID = request.Id
		desiredState.State = DeviceStateMessage{
			Version: uint64(0),
			State:   json.RawMessage([]byte("{}")),
		}
	}

	u := &jsonpb.Unmarshaler{}

	var reportedValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(reportedState.State.State), &reportedValue); err != nil {
		log.Error("Failed to unmarshal reported JSON from database", zap.Error(err))
	} else {
		response.Shadow.Reported = &shadowpb.VersionedValue{
			Version:   uint64(reportedState.State.Version),
			Data:      &reportedValue,
			Timestamp: timestamppb.New(reportedState.State.Timestamp),
		}
	}

	var desiredValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(desiredState.State.State), &desiredValue); err != nil {
		log.Error("Failed to unmarshal desired JSON from database", zap.Error(err))
	} else {
		response.Shadow.Desired = &shadowpb.VersionedValue{
			Version:   uint64(desiredState.State.Version),
			Data:      &desiredValue,
			Timestamp: timestamppb.New(desiredState.State.Timestamp),
		}
	}

	return response, nil
}

//PatchDesiredState is a method to patch a message to a device state
func (s *Server) PatchDesiredState(context context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {

	log := s.Log.Named("Patch Desired State Controller")
	log.Debug("Function Invoked", zap.String("Device", request.Id))

	// TODO sanity-check request

	var marshaler jsonpb.Marshaler
	var b bytes.Buffer
	err = marshaler.Marshal(&b, request.GetData())
	if err != nil {
		return nil, err
	}

	_, _, err = s.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: s.ProduceTopic,
		Key:   sarama.StringEncoder(request.GetId()),
		Value: sarama.ByteEncoder(b.Bytes()),
	})
	if err != nil {
		return nil, err
	}

	return &shadowpb.PatchDesiredStateResponse{}, nil
}

//StreamReportedStateChanges is a method to start streaming of data from a device
func (s *Server) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv shadowpb.Shadows_StreamReportedStateChangesServer) (err error) {

	log := s.Log.Named("Stream State Controller")
	log.Debug("Function Invoked", zap.Strings("devices", request.Devices), zap.Bool("delta", request.OnlyDelta))

	if len(request.Devices) == 0 {
		return errors.New("no devices specified")
	}

	var subPathReported string
	if request.OnlyDelta {
		subPathReported = "/reported/delta"
	} else {
		subPathReported = "/reported/full"
	}

	topicEvents := make([]string, len(request.Devices))
	for i, device := range request.Devices {
		topicEvents[i] = device + subPathReported
	}
	events := s.PubSub.Sub(topicEvents...)

	log.Debug("Reported Streaming Details", zap.Strings("topics", topicEvents))

	defer func() {

		go func() {
			s.PubSub.Unsub(events)
		}()

		// Drain

		for range events {

		}

		log.Info("Drained Reported Channel")
	}()

	var subPathDesired string
	if request.OnlyDelta {
		subPathDesired = "/desired/delta"
	} else {
		subPathDesired = "/desired/full"
	}

	topicEventsDesired := make([]string, len(request.Devices))
	for i, device := range request.Devices {
		topicEventsDesired[i] = device + subPathDesired
	}
	eventsDesired := s.PubSub.Sub(topicEventsDesired...)

	log.Debug("Desired Streaming Details", zap.Strings("topics", topicEventsDesired))

	defer func() {

		go func() {
			s.PubSub.Unsub(eventsDesired)
		}()

		// Drain

		for range eventsDesired {

		}

		log.Debug("Drained Desired Channel")
	}()
outer:
	for {
		log.Debug("Looping through Stream")
		select {
		case reportedEvent := <-events:
			log.Debug("Inside reported event reading", zap.Any("event", reportedEvent))
			device, value, err := toProto(reportedEvent, log)
			if err != nil {
				log.Error("Unable to Unmarshal data", zap.Error(err))
				break outer
			}

			log.Debug("Reported", zap.Any("Reported Value", value))

			err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
				Device: device,
				ReportedState: value,
			})
			if err != nil {
				log.Error("Unable to send Reported data", zap.Error(err))
				break outer
			}
		case desiredEvent := <-eventsDesired:
			log.Debug("Inside desired event reading", zap.Any("event", desiredEvent))
			device, value, err := toProto(desiredEvent, log)
			if err != nil {
				log.Error("Unable to Unmarshal data", zap.Error(err))
				break outer
			}

			log.Debug("Desired", zap.Any("Desired Value", value))

			err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
				Device: device,
				DesiredState: value,
			})
			if err != nil {
				log.Error("Unable to send Desired data", zap.Error(err))
				break outer
			}
		case <-srv.Context().Done():
			log.Debug("Streaming is Done")
			break outer

		}

	}
	return nil
}

func toProto(event interface{}, log *zap.Logger) (device string, result *shadowpb.VersionedValue, err error) {
	var value structpb.Value
	if raw, ok := event.(*DeviceStateMessage); ok {
		var u jsonpb.Unmarshaler
		err = u.Unmarshal(bytes.NewReader(raw.State), &value)
		if err != nil {
			log.Error("Failed to unmarshal jsonpb", zap.Error(err))
			return "", nil, err
		}

		return raw.Device, &shadowpb.VersionedValue{
			Version:   raw.Version,
			Data:      &value,
			Timestamp: timestamppb.New(raw.Timestamp), // TODO
		}, nil
	}
	return "", nil, errors.New("failed type assertion")

}
