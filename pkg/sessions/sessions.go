package sessions

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/infinimesh/proto/node/sessions"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionsHandler interface {
	New(exp int64, client string) *sessions.Session
	Store(account string, session *sessions.Session) error
	Check(account, sid string) error
	LogActivity(account, sid string, exp int64) error
	Get(account string) ([]*sessions.Session, error)
	GetActivity(account string) (map[string]*timestamppb.Timestamp, error)
	Revoke(account, sid string) error
}

type sessionsHandler struct {
	rdb redis.Cmdable

	protoMarshal   func(m protoreflect.ProtoMessage) ([]byte, error)
	protoUnmarshal func(b []byte, m protoreflect.ProtoMessage) error
}

func NewSessionsHandler(
	rdb redis.Cmdable,
	protoMarshal func(m protoreflect.ProtoMessage) ([]byte, error),
	protoUnmarshal func(b []byte, m protoreflect.ProtoMessage) error,
) SessionsHandler {
	return &sessionsHandler{
		rdb:            rdb,
		protoMarshal:   protoMarshal,
		protoUnmarshal: protoUnmarshal,
	}
}

func (s *sessionsHandler) New(_exp int64, client string) *sessions.Session {
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

func (s *sessionsHandler) Store(account string, session *sessions.Session) error {
	data, err := s.protoMarshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("sessions:%s:%s", account, session.Id)
	var ret time.Duration = 0
	if session.Expires != nil {
		ret = time.Until(session.Expires.AsTime())
	}

	return s.rdb.Set(context.Background(), key, data, ret).Err()
}

func (s *sessionsHandler) Check(account, sid string) error {
	key := fmt.Sprintf("sessions:%s:%s", account, sid)

	cmd := s.rdb.Get(context.Background(), key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	data, err := cmd.Bytes()
	if err != nil {
		return err
	}

	session := &sessions.Session{}
	err = s.protoUnmarshal(data, session)
	if err != nil {
		return err
	}

	if session.Expires != nil && session.Expires.AsTime().Before(time.Now()) {
		return fmt.Errorf("session expired")
	}

	return nil
}

func (s *sessionsHandler) LogActivity(account, sid string, exp int64) error {
	return s.rdb.Set(context.Background(), fmt.Sprintf("sessions:activity:%s:%s", account, sid), time.Now().Unix(), time.Until(time.Unix(exp, 0))).Err()
}
func (s *sessionsHandler) GetActivity(account string) (map[string]*timestamppb.Timestamp, error) {
	keys, err := s.rdb.Keys(context.Background(), fmt.Sprintf("sessions:activity:%s:*", account)).Result()
	if err != nil {
		return nil, err
	}

	data, err := s.rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*timestamppb.Timestamp)
	for i, d := range data {
		d_s, ok := d.(string)
		if !ok {
			return nil, fmt.Errorf("invalid data type: %s", keys[i])
		}

		ts, err := strconv.Atoi(d_s)
		if err != nil {
			return nil, fmt.Errorf("invalid data type: %s | %v", keys[i], err)
		}

		result[strings.Split(keys[i], ":")[3]] = timestamppb.New(time.Unix(int64(ts), 0))
	}

	return result, nil
}

func (s *sessionsHandler) Get(account string) ([]*sessions.Session, error) {

	keys, err := s.rdb.Keys(context.Background(), fmt.Sprintf("sessions:%s:*", account)).Result()
	if err != nil {
		return nil, err
	}

	data, err := s.rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*sessions.Session, len(data))
	for i, d := range data {
		session := &sessions.Session{}

		str, ok := d.(string)
		if !ok {
			return nil, fmt.Errorf("invalid data type: %s", keys[i])
		}

		err = s.protoUnmarshal([]byte(str), session)
		if err != nil {
			return nil, err
		}

		result[i] = session
	}

	return result, nil
}

func (s *sessionsHandler) Revoke(account, sid string) error {
	key := fmt.Sprintf("sessions:%s:%s", account, sid)
	return s.rdb.Del(context.Background(), key).Err()
}
