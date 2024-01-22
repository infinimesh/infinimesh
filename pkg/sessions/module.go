package sessions

import (
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type SessionsHandlerModule interface {
	Handler() SessionsHandler
}

type sessionsHandlerModule struct {
	handler SessionsHandler
}

func NewSessionsHandlerModule(rdb redis.Cmdable) SessionsHandlerModule {
	return &sessionsHandlerModule{
		handler: NewSessionsHandler(
			rdb,
			proto.Marshal,
			proto.Unmarshal,
		),
	}
}

func (s *sessionsHandlerModule) Handler() SessionsHandler {
	return s.handler
}
