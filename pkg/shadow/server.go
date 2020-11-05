//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
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
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

//Server is a data strcuture for shadow server
type Server struct {
	Repo         Repo
	Producer     sarama.SyncProducer // Sync producer, we want to guarantee execution
	ProduceTopic string
	Log          *zap.Logger

	PubSub *pubsub.PubSub
}

//Get is a method to get a device state
func (s *Server) Get(context context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {

	log := s.Log.Named("Get State Controller")
	log.Info("Function Invoked", zap.String("Device", request.Id))

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
		fmt.Fprintf(os.Stderr, "Failed to unmarshal reported JSON from database: %v\n", err)
	} else {
		ts, err := ptypes.TimestampProto(reportedState.State.Timestamp)
		if err != nil {
			return nil, err
		}
		response.Shadow.Reported = &shadowpb.VersionedValue{
			Version:   uint64(reportedState.State.Version),
			Data:      &reportedValue,
			Timestamp: ts,
		}
	}

	var desiredValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(desiredState.State.State), &desiredValue); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON from database: %v\n", err)
	} else {
		ts, err := ptypes.TimestampProto(desiredState.State.Timestamp)
		if err != nil {
			return nil, err
		}

		response.Shadow.Desired = &shadowpb.VersionedValue{
			Version:   uint64(desiredState.State.Version),
			Data:      &desiredValue,
			Timestamp: ts,
		}
	}

	return response, nil
}

//PatchDesiredState is a method to patch a message to a device state
func (s *Server) PatchDesiredState(context context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {

	log := s.Log.Named("Patch Desired State Controller")
	log.Info("Function Invoked", zap.String("Device", request.Id))

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
	log.Info("Function Invoked", zap.String("Device", request.Id), zap.Bool("Delta Flag", request.OnlyDelta))

	// TODO validate request/Id

	var subPathReported string
	if request.OnlyDelta {
		subPathReported = "/state/reported/delta"
	} else {
		subPathReported = "/state/reported/full"
	}
	log.Info("streaming requested by : ", zap.String("Request ID :", request.Id))
	topicEvents := "devices/" + request.Id + subPathReported
	events := s.PubSub.Sub(topicEvents)
	fmt.Println(topicEvents)
	fmt.Println(events)
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
		subPathDesired = "/state/desired/delta"
	} else {
		subPathDesired = "/state/desired/full"
	}

	topicEventsDesired := "devices/" + request.Id + subPathDesired
	eventsDesired := s.PubSub.Sub(topicEventsDesired)
	defer func() {

		go func() {
			s.PubSub.Unsub(eventsDesired)
		}()

		// Drain

		for range eventsDesired {

		}

		log.Info("Drained Desired Channel")
	}()
outer:
	for {

		select {
		case reportedEvent := <-events:
			value, err := toProto(reportedEvent)
			if err != nil {
				fmt.Println(err)
				break outer
			}
			fmt.Printf("Server Reported Value : %v", value)
			err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
				ReportedState: value,
			})
			if err != nil {
				fmt.Println(err)
				break outer
			}
		case desiredEvent := <-eventsDesired:
			value, err := toProto(desiredEvent)
			if err != nil {
				fmt.Println(err)
				break outer
			}
			fmt.Printf("Server Desired Value : %v", value)
			err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
				DesiredState: value,
			})
			if err != nil {
				fmt.Println(err)
				break outer
			}
		case <-srv.Context().Done():
			fmt.Println("DONE")
			break outer

		}

	}
	return nil
}

func toProto(event interface{}) (result *shadowpb.VersionedValue, err error) {
	var value structpb.Value
	if raw, ok := event.(*DeviceStateMessage); ok {
		var u jsonpb.Unmarshaler
		err = u.Unmarshal(bytes.NewReader(raw.State), &value)
		if err != nil {
			fmt.Println("Failed to unmarshal jsonpb: ", err)
			return nil, err
		}

		ts, err := ptypes.TimestampProto(raw.Timestamp)
		if err != nil {
			fmt.Println("Invalid timestamp", err)
			return nil, err
		}

		return &shadowpb.VersionedValue{
			Version:   raw.Version,
			Data:      &value,
			Timestamp: ts, // TODO
		}, nil
	}
	return nil, errors.New("Failed type assertion")

}
