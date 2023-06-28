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

func Check(rdb *redis.Client, account, sid string) error {
	key := fmt.Sprintf("sessions:%s:%s", account, sid)

	cmd := rdb.Get(context.Background(), key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	data, err := cmd.Bytes()
	if err != nil {
		return err
	}

	session := &sessions.Session{}
	err = proto.Unmarshal(data, session)

	if err != nil {
		return err
	}

	if session.Expires != nil && session.Expires.AsTime().Before(time.Now()) {
		return fmt.Errorf("session expired")
	}

	return nil
}

func LogActivity(rdb *redis.Client, account, sid string, exp int64) error {
	return rdb.Set(context.Background(), fmt.Sprintf("sessions:activity:%s:%s", account, sid), time.Now().Unix(), time.Until(time.Unix(exp, 0))).Err()
}

func Get(rdb *redis.Client, account string) ([]*sessions.Session, error) {

	keys, err := rdb.Keys(context.Background(), fmt.Sprintf("sessions:%s:*", account)).Result()
	if err != nil {
		return nil, err
	}

	data, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*sessions.Session, len(data))
	for i, d := range data {
		session := &sessions.Session{}

		bytes, ok := d.(string)
		if !ok {
			return nil, fmt.Errorf("invalid data type: %s", keys[i])
		}

		err = proto.Unmarshal([]byte(bytes), session)
		if err != nil {
			return nil, err
		}

		result[i] = session
	}

	return result, nil
}
