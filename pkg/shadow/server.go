package shadow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	// TODO fetch device from registry, 404 if not found

	reportedState, err := s.Repo.GetReported(req.Id)
	if err != nil {
		reportedState.ID = req.Id
		reportedState.State = json.RawMessage([]byte("{}"))
		reportedState.Version = uint64(0)
	}

	desiredState, err := s.Repo.GetDesired(req.Id)
	if err != nil {
		desiredState.ID = req.Id
		desiredState.State = json.RawMessage([]byte("{}"))
		desiredState.Version = uint64(0)
	}

	u := &jsonpb.Unmarshaler{}

	var reportedValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(reportedState.State), &reportedValue); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal reported JSON from database: %v\n", err)
	} else {
		response.Shadow.Reported = &shadowpb.VersionedValue{
			Version:   uint64(reportedState.Version),
			Data:      &reportedValue,
			Timestamp: ts,
		}
	}

	var desiredValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(desiredState.State), &desiredValue); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal JSON from database: %v\n", err)
	} else {
		response.Shadow.Desired = &shadowpb.VersionedValue{
			Version:   uint64(desiredState.Version),
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
		if raw, ok := event.(json.RawMessage); ok {
			var u jsonpb.Unmarshaler
			err = u.Unmarshal(bytes.NewReader(raw), &value)
			if err != nil {
				fmt.Println("Failed to unmarshal jsonpb: ", err)
				continue
			}
		} else {
			fmt.Println("Failed type assertion")
			continue
		}

		err = srv.Send(&shadowpb.StreamReportedStateChangesResponse{
			ReportedDelta: &shadowpb.VersionedValue{
				Version:   0,
				Data:      &value,
				Timestamp: ptypes.TimestampNow(), // TODO
			},
		})
		if err != nil {
			break
		}
	}
	return nil
}
