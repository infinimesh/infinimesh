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
