package graph

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/infinimesh/infinimesh/pkg/sessions"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	node "github.com/infinimesh/proto/node"
	pb "github.com/infinimesh/proto/node/sessions"
	"go.uber.org/zap"
)

type SessionsController struct {
	node.UnimplementedSessionsServiceServer

	log *zap.Logger
	rdb *redis.Client
}

func NewSessionsController(log *zap.Logger, rdb *redis.Client) *SessionsController {
	return &SessionsController{
		log: log.Named("Sessions"),
		rdb: rdb,
	}
}

func (c *SessionsController) Get(ctx context.Context, req *node.EmptyMessage) (*pb.Sessions, error) {
	log := c.log.Named("Get")
	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	sid := ctx.Value(inf.InfinimeshSessionCtxKey).(string)

	log.Debug("Invoked", zap.String("requestor", requestor), zap.String("sid", sid))

	result, err := sessions.Get(c.rdb, requestor)
	if err != nil {
		return nil, err
	}

	current := true
	for _, session := range result {
		if session.Id == sid {
			session.Current = &current
			break
		}
	}

	return &pb.Sessions{
		Sessions: result,
	}, nil
}

func (c *SessionsController) GetActivity(ctx context.Context, req *node.EmptyMessage) (*pb.Activity, error) {
	log := c.log.Named("GetActivity")
	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)

	log.Debug("Invoked", zap.String("requestor", requestor))

	result, err := sessions.GetActivity(c.rdb, requestor)
	if err != nil {
		return nil, err
	}

	return &pb.Activity{
		LastSeen: result,
	}, nil
}

func (c *SessionsController) Revoke(ctx context.Context, req *pb.Session) (*node.DeleteResponse, error) {
	log := c.log.Named("Revoke")
	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)

	log.Debug("Invoked", zap.String("requestor", requestor), zap.String("sid", req.Id))

	err := sessions.Revoke(c.rdb, requestor, req.Id)
	if err != nil {
		return nil, err
	}

	return &node.DeleteResponse{}, nil
}
