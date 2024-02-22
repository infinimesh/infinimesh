/*
Copyright © 2021-2023 Infinite Devices GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package shadow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	pb "github.com/infinimesh/proto/shadow"
	"go.uber.org/zap"
)

func Key(device string, key pb.StateKey) string {
	var k string

	switch key {
	case pb.StateKey_DESIRED:
		k = "desired"
	case pb.StateKey_REPORTED:
		k = "reported"
	case pb.StateKey_CONNECTION:
		k = "connection"
	default:
		k = "garbage"
	}

	return fmt.Sprintf("%s:%s", device, k)
}

func (s *ShadowServiceServer) Persister() {
	log := s.log.Named("persister")

	log.Info("Starting Persister")
	defer log.Warn("Exited")

	messages := make(chan interface{}, 10)
	s.ps.AddSub(messages, "mqtt.incoming", "mqtt.outgoing")
	defer unsub(s.ps, messages)

	for msg := range messages {
		shadow := msg.(*pb.Shadow)
		log.Debug("Message received", zap.Any("shadow", shadow))
		if shadow.Reported != nil {
			log.Debug("Reporting", zap.String("device", shadow.Device))
			s.MergeAndStore(log, shadow.Device, pb.StateKey_REPORTED, shadow.Reported)
		}
		if shadow.Desired != nil {
			log.Debug("Desiring", zap.String("device", shadow.Device))
			s.MergeAndStore(log, shadow.Device, pb.StateKey_DESIRED, shadow.Desired)
		}
		if shadow.Connection != nil {
			s.StoreConnectionState(log, shadow.Device, shadow.Connection)
		}
	}
}

func (s *ShadowServiceServer) MergeAndStore(log *zap.Logger, device string, skey pb.StateKey, state *pb.State) {
	key := Key(device, skey)

	var new []byte
	var err error

	new, err = json.Marshal(state)
	if err != nil {
		log.Warn("Error Marshalling State", zap.String("key", key), zap.Error(err))
		return
	}

	cmd := s.rdb.Get(context.Background(), key)
	m, err := cmd.Result()
	if err != nil {
		goto set
	}

	log.Debug("Merging", zap.ByteString("old", []byte(m)), zap.ByteString("new", new))
	new, err = MergeJSON([]byte(m), new)
	if err != nil {
		log.Warn("Error Merging State", zap.String("key", key), zap.Error(err))
		return
	}

set:
	r := s.rdb.Set(context.Background(), key, string(new), 0)
	if r.Err() != nil {
		log.Warn("Error Storing State", zap.String("key", key), zap.Error(err))
		return
	}

	err = s.repo.UpdateDeviceModifyDate(context.Background(), log, device)
	if err != nil {
		log.Warn("Error updating modify date", zap.Error(err))
		return
	}

}

func MergeJSON(old, new []byte) ([]byte, error) {
	merged, err := jsonpatch.MergePatch(old, new)
	if err != nil {
		if string(old) == "" {
			return new, nil
		}
		if string(new) == "" {
			return old, nil
		}
		return nil, err
	}
	return merged, nil
}

func (s *ShadowServiceServer) Store(log *zap.Logger, device string, skey pb.StateKey, state interface{}) (string, bool) {
	key := Key(device, skey)

	new, err := json.Marshal(state)
	if err != nil {
		log.Warn("Error Marshalling State", zap.String("key", key), zap.Error(err))
		return key, false
	}

	r := s.rdb.Set(context.Background(), key, string(new), 0)
	if r.Err() != nil {
		log.Warn("Error Storing State", zap.String("key", key), zap.Error(r.Err()))
		return key, false
	}

	err = s.repo.UpdateDeviceModifyDate(context.Background(), log, device)
	if err != nil {
		log.Warn("Error updating modify date", zap.Error(err))
		return key, false
	}

	return key, true
}

func (s *ShadowServiceServer) StoreConnectionState(log *zap.Logger, device string, state *pb.ConnectionState) {

	key, ok := s.Store(log, device, pb.StateKey_CONNECTION, state)
	if !ok {
		return
	}

	r := s.rdb.Expire(context.Background(), key, time.Hour*24)
	if r.Err() != nil {
		log.Warn("Couldn't set key expiration", zap.String("key", key), zap.Error(r.Err()))
	}
}
