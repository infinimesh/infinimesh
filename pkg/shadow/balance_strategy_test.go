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

package shadow

import (
	"testing"

	sarama "github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
)

func TestBalance(t *testing.T) {
	strat := &BalanceStrategyCoPartitioned{}

	members := map[string]sarama.ConsumerGroupMemberMetadata{
		"member-a": sarama.ConsumerGroupMemberMetadata{},
	}

	topics := map[string][]int32{
		"first-topic": []int32{1, 2},
	}

	plan, err := strat.Plan(members, topics)
	require.NoError(t, err)

	require.Len(t, plan, 1)                            // Number of members
	require.Len(t, plan["member-a"], 1)                // Number of topics
	require.Len(t, plan["member-a"]["first-topic"], 2) // Number of partitions
}

func TestBalanceMultipleMembers(t *testing.T) {
	strat := &BalanceStrategyCoPartitioned{}

	members := map[string]sarama.ConsumerGroupMemberMetadata{
		"member-a": sarama.ConsumerGroupMemberMetadata{},
		"member-b": sarama.ConsumerGroupMemberMetadata{},
	}

	topics := map[string][]int32{
		"first-topic":  []int32{1, 2},
		"second-topic": []int32{1, 2},
	}

	plan, err := strat.Plan(members, topics)

	require.NoError(t, err)
	require.Len(t, plan, 2)                                                                   // Number of members
	require.Len(t, plan["member-a"], 2)                                                       // Number of topics
	require.Len(t, plan["member-b"], 2)                                                       // Number of topics
	require.Len(t, plan["member-a"]["first-topic"], 1)                                        // Number of partitions
	require.EqualValues(t, plan["member-a"]["first-topic"], plan["member-a"]["second-topic"]) // Same partitions on both topics
}
