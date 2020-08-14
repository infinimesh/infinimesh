//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package timeseries

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"encoding/json"

	"github.com/jeremywohl/flatten"

	"github.com/infinimesh/infinimesh/pkg/shadow"
)

type Consumer struct {
	Log  *zap.Logger
	Repo TimeseriesRepo
}

func (h *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	return nil
}

func (h *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if message == nil {
			break
		}

		var msg shadow.DeviceStateMessage

		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			h.Log.Error("Failed to unmarshal message", zap.Error(err), zap.Int64("offset", message.Offset), zap.Int32("partition", message.Partition))
		}

		fmt.Println("got msg", string(message.Value))

		flatJSON, err := flatten.FlattenString(string(msg.State), "", flatten.DotStyle)
		if err != nil {
			h.Log.Info("Failed to flatten", zap.Error(err))
		}

		flat := make(map[string]interface{})
		err = json.Unmarshal([]byte(flatJSON), &flat)
		if err != nil {
			h.Log.Error("Failed to unmarshal", zap.Error(err))
		}

		for property, value := range flat {
			var datapointValue float64

			switch v := value.(type) {
			case float64:
				datapointValue = v
			case bool:
				if v {
					datapointValue = 1
				} else {
					datapointValue = 0
				}
			}

			err = h.Repo.CreateDataPoint(context.TODO(), &DataPoint{
				DeviceID:   string(message.Key),
				DeviceName: "hardcoded-test",
				Property:   property,
				Timestamp:  msg.Timestamp,
				Value:      datapointValue,
			})

		}

		if err != nil {
			h.Log.Error("Failed to store datapoint", zap.Error(err))
		}

		s.MarkMessage(message, "")
	}
	return nil
}
