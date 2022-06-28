/*
Copyright © 2021-2022 Infinite Devices GmbH

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

	jsonpatch "github.com/evanphx/json-patch"
	pb "github.com/infinimesh/proto/shadow"
	"go.uber.org/zap"
)

func Key(device string, key pb.StateKey) string {
	var k string
	if key == pb.StateKey_DESIRED {
		k = "desired"
	} else if key == pb.StateKey_REPORTED {
		k = "reported"
	} else {
		k = "garbage"
	}
	return fmt.Sprintf("%s:%s", device, k)
}

func (s *ShadowServiceServer) Persister() {
	log := s.log.Named("persister")
	messages := make(chan interface{}, 10)
	s.ps.AddSub(messages, "mqtt.incoming", "mqtt.outgoing")

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
	}
}

func (s *ShadowServiceServer) MergeAndStore(log *zap.Logger, device string, skey pb.StateKey, state *pb.State) {
	key := Key(device, skey)

	var new, merged []byte
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

merge:
	log.Debug("Merging", zap.ByteString("old", []byte(m)), zap.ByteString("new", new))
	merged, err = jsonpatch.MergePatch([]byte(m), new)
	if err != nil {
		if m == "" {
			m = "{}"
			goto merge
		}
		log.Warn("Error Merging State", zap.String("key", key), zap.Error(err))
		return
	}
	new = merged

set:
	r := s.rdb.Set(context.Background(), key, string(new), 0)
	if r.Err() != nil {
		log.Warn("Error Storing State", zap.String("key", key), zap.Error(err))
		return
	}
}

func (s *ShadowServiceServer) Store(log *zap.Logger, device string, skey pb.StateKey, state *pb.State) {
	key := Key(device, skey)

	new, err := json.Marshal(state)
	if err != nil {
		log.Warn("Error Marshalling State", zap.String("key", key), zap.Error(err))
		return
	}

	r := s.rdb.Set(context.Background(), key, string(new), 0)
	if r.Err() != nil {
		log.Warn("Error Storing State", zap.String("key", key), zap.Error(err))
		return
	}
}
