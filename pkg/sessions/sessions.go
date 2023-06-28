package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/infinimesh/proto/node/sessions"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(_exp int64, client string) *sessions.Session {
	now := time.Now()
	id := fmt.Sprintf("%x", now.UnixNano())

	exp := timestamppb.New(time.Unix(_exp, 0))
	if _exp == 0 {
		exp = nil
	}

	return &sessions.Session{
		Id:      id,
		Expires: exp,
		Client:  client,
		Created: timestamppb.New(now),
	}
}

func Store(rdb *redis.Client, account string, session *sessions.Session) error {
	data, err := proto.Marshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("sessions:%s:%s", account, session.Id)
	var ret time.Duration = 0
	if session.Expires != nil {
		ret = time.Until(session.Expires.AsTime())
	}

	return rdb.Set(context.Background(), key, data, ret).Err()
}
