package zkOracle

import (
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type StateSync struct {
	state    *State
	contract *ZKOracleContract
}

func NewStateSync(state *State, contract *ZKOracleContract) *StateSync {
	return &StateSync{state: state, contract: contract}
}

func (s *StateSync) Synchronize() error {
	go func() {
		if err := WatchEvent(context.Background(), s.contract.WatchRegistered, s.HandleRegisteredEvent); err != nil {
			logger.Err(err).Msg("watch registered event")
		}
	}()

	return nil
}

func (s *StateSync) HandleRegisteredEvent(ctx context.Context, event *ZKOracleContractRegistered) error {

	logger.Info().
		Uint64("index", event.Index.Uint64()).
		Str("pubKeyX", event.Pubkey.X.String()).
		Str("pubKeyY", event.Pubkey.Y.String()).
		Str("balance", event.Value.String()).
		Msg("handle registered event")

	x := fr.NewElement(0)
	y := fr.NewElement(0)

	x.SetBigInt(event.Pubkey.X)
	y.SetBigInt(event.Pubkey.Y)

	publicKey := eddsa.PublicKey{A: twistededwards.NewPointAffine(x, y)}

	account := Account{
		Index:     event.Index,
		PublicKey: &publicKey,
		Balance:   event.Value,
	}

	err := s.state.WriteAccount(account)
	if err != nil {
		return fmt.Errorf("write account: %w", err)
	}

	return nil
}
