package zkOracle

import (
	"context"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type Aggregator struct {
	votePool         *VotePool
	constraintSystem frontend.CompiledConstraintSystem
	provingKey       groth16.ProvingKey
	verifyingKey     groth16.VerifyingKey
}

func NewAggregator(votePool *VotePool, constraintSystem frontend.CompiledConstraintSystem, provingKey groth16.ProvingKey, verifyingKey groth16.VerifyingKey) *Aggregator {
	return &Aggregator{votePool: votePool, constraintSystem: constraintSystem, provingKey: provingKey, verifyingKey: verifyingKey}
}

func (a *Aggregator) Aggregate(ctx context.Context) error {
	for requestNumber := range a.votePool.sink {
		logger.Info().Uint64("requestNumber", requestNumber).Msg("aggregate votes")
	}

	return nil
}
