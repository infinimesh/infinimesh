package shadow

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type Server struct {
	Repo Repo
}

func (s *Server) Get(context context.Context, req *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	reportedState, err := s.Repo.GetReported(req.Id)
	if err != nil {
		return nil, err
	}

	desiredState, err := s.Repo.GetDesired(req.Id)
	if err != nil {
		return nil, err
	}

	u := &jsonpb.Unmarshaler{}

	var reportedValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(reportedState.State), &reportedValue); err != nil {
		return nil, errors.New("Failed to unmarshal JSON from database")
	}

	var desiredValue structpb.Value
	if err := u.Unmarshal(bytes.NewReader(desiredState.State), &desiredValue); err != nil {
		return nil, errors.New("Failed to unmarshal JSON from database")
	}

	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	return &shadowpb.GetResponse{
		Shadow: &shadowpb.Shadow{
			Reported: &shadowpb.VersionedValue{
				Version:   uint64(reportedState.Version),
				Data:      &reportedValue,
				Timestamp: ts,
			},
			Desired: &shadowpb.VersionedValue{
				Version:   uint64(desiredState.Version),
				Data:      &desiredValue,
				Timestamp: ts,
			},
		},
	}, nil
}
