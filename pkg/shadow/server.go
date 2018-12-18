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
	state, err := s.Repo.GetReported(req.Id)
	if err != nil {
		return nil, err
	}

	u := &jsonpb.Unmarshaler{}

	var value structpb.Value
	if err := u.Unmarshal(bytes.NewReader(state.State), &value); err != nil {
		return nil, errors.New("Failed to unmarshal JSON from database")
	}

	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	return &shadowpb.GetResponse{
		Shadow: &shadowpb.Shadow{
			Reported: &shadowpb.VersionedValue{
				Version:   uint64(state.Version),
				Data:      &value,
				Timestamp: ts,
			},
		},
	}, nil
}
