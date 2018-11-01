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
