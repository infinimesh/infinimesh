package shadow

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type Server struct {
	Repo Repo
}

func (s *Server) GetReported(context context.Context, req *shadowpb.GetReportedRequest) (response *shadowpb.GetReportedResponse, err error) {
	state, err := s.Repo.GetReported(req.DeviceId)
	if err != nil {
		return nil, err
	}

	u := &jsonpb.Unmarshaler{}
	fmt.Println("test", string(state.State))

	var value structpb.Value
	if err := u.Unmarshal(bytes.NewReader(state.State), &value); err != nil {
		return nil, errors.New("Failed to unmarshal JSON from database")
	}

	return &shadowpb.GetReportedResponse{State: &value}, nil
}
