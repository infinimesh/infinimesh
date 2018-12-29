package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/dgo"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/infinimesh/infinimesh/pkg/auth/authpb"
)

type Server struct {
	Dgraph *dgo.Dgraph
}

func (s *Server) Authorize(ctx context.Context, request *authpb.AuthorizeRequest) (response *authpb.AuthorizeResponse, err error) {
	params := map[string]string{
		"$device_id":  request.GetObject(),
		"$user_email": request.GetSubject(),
		"$action":     request.GetAction(),
	}
	fmt.Println(params)
	const q = `{
  var(func: eq(device_id,"testdevice4")) @recurse @normalize @cascade {
    parentObjectUIDs as uid
    contained_in  {
    }
  }

  var(func: uid(parentObjectUIDs)) @normalize  @cascade {
    clearances @filter(eq(action, "write")) {
      clearanceIDs as uid
      
    }
  }
      
  firstWriteClearance(func: uid(clearanceIDs), first: 1) @cascade {
    uid
    action
    granted_to @filter(eq(email, "birdy@nerden.de")) {}
  }
}`

	res, err := s.Dgraph.NewTxn().Query(ctx, q)

	type Permission struct {
		Action string `json:"action"`
	}

	type Permissions struct {
		Permissions []Permission `json:"firstWriteClearance"`
	}

	if err != nil {
		return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	var p Permissions
	err = json.Unmarshal(res.Json, &p)
	if err != nil {
		return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	spew.Dump(p)

	return nil, nil
}
