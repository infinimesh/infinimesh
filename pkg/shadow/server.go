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
		reportedState.State = DeviceStateMessage{
			Version: uint64(0),
			State:   json.RawMessage([]byte("{}")),
		}
	}

	desiredState, err := s.Repo.GetDesired(req.Id)
	if err != nil {
		desiredState.ID = req.Id
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

	var subPathReported string
	if request.OnlyDelta {
		subPathReported = "/reported/delta"
	} else {
		subPathReported = "/reported/full"
	}

	topicEvents := request.Id + subPathReported
	events := s.PubSub.Sub(topicEvents)
	defer func() {
		fmt.Println("Dferer")
		go func() {
			s.PubSub.Unsub(events)
		}()

		// Drain

		for range events {

		}

		fmt.Println("Drained channel")
	}()

	var subPathDesired string
	if request.OnlyDelta {
		subPathDesired = "/desired/delta"
	} else {
		subPathDesired = "/desired/full"
	}

	topicEventsDesired := request.Id + subPathDesired
	eventsDesired := s.PubSub.Sub(topicEventsDesired)
	defer func() {
		fmt.Println("defer2")
		go func() {
			s.PubSub.Unsub(eventsDesired)
		}()

		// Drain

		for range eventsDesired {

		}

		fmt.Println("Drained channel")
	}()
outer:
	for {

		fmt.Println("Vor select")
		select {
		case reportedEvent := <-events:
			value, err := toProto(reportedEvent)
			if err != nil {
				fmt.Println(err)
				break outer
			}

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
	} else {
		return nil, errors.New("Failed type assertion")
	}

}
