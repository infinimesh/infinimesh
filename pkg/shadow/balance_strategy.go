//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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

package shadow

import (
	"errors"

	"github.com/Shopify/sarama"
)

type BalanceStrategyCoPartitioned struct{}

func (strat *BalanceStrategyCoPartitioned) Name() string {
	return "parallel"
}

func (strat *BalanceStrategyCoPartitioned) Plan(members map[string]sarama.ConsumerGroupMemberMetadata, topics map[string][]int32) (sarama.BalanceStrategyPlan, error) {

	// Check if topics are co-partitioned

	var previousLength int
	for _, partitions := range topics {
		if previousLength != 0 {
			// Compare
			if len(partitions) != previousLength {
				return nil, errors.New("Topics are not co-partitioned")
			}
		}

		previousLength = len(partitions)
	}

	plan := sarama.BalanceStrategyPlan{}

	var memberIDs []string
	for member := range members {
		memberIDs = append(memberIDs, member)
	}

	for i := 0; i < previousLength; i++ {
		member := memberIDs[i%len(members)]
		for topic := range topics {
			plan.Add(member, topic, int32(i))
		}
	}

	return plan, nil
}
