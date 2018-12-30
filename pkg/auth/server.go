package auth

import (
	"context"
	"encoding/json"
	"fmt"

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
	const q = `query permissions($action: string, $device_id: string, $user_email: string){
                     var(func: eq(device_id,$device_id)) @recurse @normalize @cascade {
                       parentObjectUIDs as uid
                       contained_in  {
                       }
                     }

                     var(func: uid(parentObjectUIDs)) @normalize  @cascade {
                       clearances @filter(eq(action, $action)) {
                         clearanceIDs as uid
                       }
                     }

                     clearance(func: uid(clearanceIDs), first: 1) @cascade {
                       uid
                       action
                       granted_to @filter(eq(email, $user_email)) {}
                     }
                   }`

	res, err := s.Dgraph.NewTxn().QueryWithVars(ctx, q, params)

	type Clearance struct {
		Action string `json:"action"`
	}

	type Permissions struct {
		Permissions []Clearance `json:"clearance"`
	}

	if err != nil {
		return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	var p Permissions
	err = json.Unmarshal(res.Json, &p)
	if err != nil {
		return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	if len(p.Permissions) > 0 {
		permission := p.Permissions[0]
		if permission.Action == request.Action {
			fmt.Printf("Granting access on object %v to subject %v for action %v. Clearance: %v\n", request.GetObject(), request.GetSubject(), request.GetAction(), permission.Action)
			return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: true}}, err
		}
	}

	return &authpb.AuthorizeResponse{Decision: &wrappers.BoolValue{Value: false}}, err
}
