package zkOracle

import (
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type StateSync struct {
	state     *State
	contract  *ZKOracleContract
	ethClient *ethclient.Client
}

func NewStateSync(state *State, contract *ZKOracleContract, ethClient *ethclient.Client) *StateSync {
	return &StateSync{state: state, contract: contract, ethClient: ethClient}
}

func (s *StateSync) Synchronize() error {

	if err := s.Update(context.Background()); err != nil {
		return fmt.Errorf("update state: %w", err)
	}

	go func() {
		if err := WatchEvent(context.Background(), s.contract.WatchRegistered, s.HandleRegisteredEvent); err != nil {
			logger.Err(err).Msg("watch registered event")
		}
	}()

	return nil
}

func (s *StateSync) Update(ctx context.Context) error {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(50),
		Addresses: []common.Address{
			common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"),
		},
	}

	logs, err := s.ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("filter logs: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(ZKOracleContractABI))
	if err != nil {
		return fmt.Errorf("reading contract abi: %w", err)
	}

	for _, log := range logs {
		event, err := contractABI.EventByID(log.Topics[0])
		if err != nil {
			return fmt.Errorf("event by id: %w", err)
		}
		switch event.Name {
		case "Registered":
			e, err := s.contract.ParseRegistered(log)
			if err != nil {
				return fmt.Errorf("parse registered: %w", err)
			}
			err = s.HandleRegisteredEvent(ctx, e)
			if err != nil {
				return fmt.Errorf("handle registered: %w", err)
			}
		}
	}

	/*	iter, err := s.contract.FilterRegistered(&bind.FilterOpts{
			Start:   0,
			End:     nil,
			Context: ctx,
		})
		if err != nil {
			return fmt.Errorf("filter registered events: %w", err)
		}

		for iter.Next() {
			if err := s.HandleRegisteredEvent(ctx, iter.Event); err != nil {
				return fmt.Errorf("%handle registered event: %w", err)
			}
		}*/
	return nil
}

func (s *StateSync) HandleRegisteredEvent(ctx context.Context, event *ZKOracleContractRegistered) error {
	logger.Info().
		Uint64("Index", event.Index.Uint64()).
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
