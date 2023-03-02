package zkOracle

import (
	"context"
	"encoding/hex"
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
	index           uint64
	state           *State
	contract        *ZKOracleContract
	contractAddress common.Address
	ethClient       *ethclient.Client
}

func NewStateSync(index uint64, state *State, contract *ZKOracleContract, contractAddress common.Address, ethClient *ethclient.Client) *StateSync {
	return &StateSync{index: index, state: state, contract: contract, contractAddress: contractAddress, ethClient: ethClient}
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

	go func() {
		if err := WatchEvent(context.Background(), s.contract.WatchBlockSubmitted, s.HandleBlockSubmittedEvent); err != nil {
			logger.Err(err).Msg("watch block submitted event")
		}
	}()

	return nil
}

func (s *StateSync) Update(ctx context.Context) error {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(50),
		Addresses: []common.Address{
			s.contractAddress,
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
		case "BlockSubmitted":
			e, err := s.contract.ParseBlockSubmitted(log)
			if err != nil {
				return fmt.Errorf("parse block submitted: %w", err)
			}
			err = s.HandleBlockSubmittedEvent(ctx, e)
			if err != nil {
				return fmt.Errorf("handle block submitted: %w", err)
			}
		}
	}

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

func (s *StateSync) HandleBlockSubmittedEvent(ctx context.Context, event *ZKOracleContractBlockSubmitted) error {
	logger.Info().
		Uint64("request", event.Request.Uint64()).
		Uint64("submitter", event.Submitter.Uint64()).
		Msg("handle block submitted event")

	if event.Submitter.Uint64() == s.index {
		return nil
	}

	aggregatorAccount, err := s.state.ReadAccount(event.Submitter.Uint64())
	if err != nil {
		return fmt.Errorf("read aggregator account: %w", err)
	}
	aggregatorAccount.Balance.Add(aggregatorAccount.Balance, big.NewInt(AggregatorReward))
	err = s.state.WriteAccount(aggregatorAccount)
	if err != nil {
		return fmt.Errorf("write account: %w", err)
	}

	for i := 0; i < NumAccounts; i++ {
		if event.Validators.Bit(i) == 0 {
			continue
		}

		account, err := s.state.ReadAccount(uint64(i))
		if err != nil {
			return fmt.Errorf("read aggregator account: %w", err)
		}
		account.Balance.Add(account.Balance, big.NewInt(ValidatorReward))
		err = s.state.WriteAccount(account)
		if err != nil {
			return fmt.Errorf("write account: %w", err)
		}
	}
	root, _ := s.state.Root()
	logger.Info().Str("root", hex.EncodeToString(root)).Msg("after update")
	return nil
}
