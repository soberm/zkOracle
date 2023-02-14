package zkOracle

import (
	"context"
)

type Aggregator struct {
	votePool *VotePool
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate(ctx context.Context) error {
	return nil
}
