package shadow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type Server struct {
	Repo         Repo
	Producer     sarama.SyncProducer // Sync producer, we want to guarantee execution
	ProduceTopic string

	PubSub *pubsub.PubSub
}

func (s *Server) Get(context context.Context, req *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	response = &shadowpb.GetResponse{
		Shadow: &shadowpb.Shadow{},
	}

	// TODO fetch device from registry, 404 if not found

	reportedState, err := s.Repo.GetReported(req.Id)
	if err != nil {
		reportedState.ID = req.Id
		reportedState.State = FullDeviceStateMessage{
			Version: uint64(0),
			State:   json.RawMessage([]byte("{}")),
		}
	}

	desiredState, err := s.Repo.GetDesired(req.Id)
	if err != nil {
		desiredState.ID = req.Id
		desiredState.State = FullDeviceStateMessage{
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

func (s *Server) PatchDesiredState(context context.Context, req *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {
	// TODO sanity-check request

	var marshaler jsonpb.Marshaler
	var b bytes.Buffer
	err = marshaler.Marshal(&b, req.GetData())
	if err != nil {
		return nil, err
	}

	_, _, err = s.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: s.ProduceTopic,
		Key:   sarama.StringEncoder(req.GetId()),
		Value: sarama.ByteEncoder(b.Bytes()),
	})
	if err != nil {
		return nil, err
	}
	return &shadowpb.PatchDesiredStateResponse{}, nil
}

func (s *Server) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv shadowpb.Shadows_StreamReportedStateChangesServer) (err error) {
	// TODO validate request/Id
	events := s.PubSub.Sub(request.Id)
	defer s.PubSub.Unsub(events)
	for event := range events {
		var value structpb.Value
		if raw, ok := event.(*DeltaDeviceStateMessage); ok {
			var u jsonpb.Unmarshaler
			err = u.Unmarshal(bytes.NewReader(raw.Delta), &value)
			if err != nil {
				fmt.Println("Failed to unmarshal jsonpb: ", err)
				continue
			}

			ts, err := ptypes.TimestampProto(raw.Timestamp)
			if err != nil {
				fmt.Println("Invalid timestamp", err)
				break
			}

			err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
				ReportedDelta: &shadowpb.VersionedValue{
					Version:   raw.Version,
					Data:      &value,
					Timestamp: ts, // TODO
				},
			})
			if err != nil {
				break
			}
		} else {
			fmt.Println("Failed type assertion")
			continue
		}

	}
	return nil
}
