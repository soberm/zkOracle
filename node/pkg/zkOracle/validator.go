package zkOracle

import (
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
)

const CONFIRMATIONS uint64 = 5

type Validator struct {
	index      uint64
	ethClient  *ethclient.Client
	contract   *ZKOracleContract
	privateKey *eddsa.PrivateKey
}

func NewValidator(index uint64, ethClient *ethclient.Client, zkOracleContract *ZKOracleContract, privateKey *eddsa.PrivateKey) *Validator {
	return &Validator{index: index, ethClient: ethClient, contract: zkOracleContract, privateKey: privateKey}
}

func (v *Validator) Validate(ctx context.Context) error {
	if err := WatchEvent(ctx, v.contract.WatchBlockRequested, v.HandleBlockRequestedEvent); err != nil {
		return fmt.Errorf("watch block requested events: %w", err)
	}

	return nil
}

func (v *Validator) HandleBlockRequestedEvent(ctx context.Context, event *ZKOracleContractBlockRequested) error {
	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("blockNumber", event.Number.Uint64()).
		Msg("handle block requested event")

	block, err := v.ethClient.HeaderByNumber(ctx, event.Number)
	if err != nil {
		return fmt.Errorf("block by number: %w", err)
	}

	currentBlockNumber, err := v.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("blocknumber: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("blockNumber", event.Number.Uint64()).
		Uint64("head", currentBlockNumber).
		Msg("check block confirmed")

	if currentBlockNumber-block.Number.Uint64() < CONFIRMATIONS {
		return fmt.Errorf("block not confirmed")
	}

	value := new(big.Int)
	value.Mod(block.Hash().Big(), fr.Modulus())

	hash := common.BytesToHash(value.Bytes())

	vote := &Vote{
		Index:     v.index,
		Request:   event.Request,
		BlockHash: hash,
	}

	hasher := mimc.NewMiMC()
	hasher.Write(vote.Serialize())
	msg := hasher.Sum(nil)

	sig, err := v.privateKey.Sign(msg, mimc.NewMiMC())
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	i, err := v.contract.GetAggregator(
		&bind.CallOpts{
			Context: ctx,
		},
	)
	if err != nil {
		return fmt.Errorf("get aggregator: %w", err)
	}

	addr, err := v.contract.GetIPAddress(&bind.CallOpts{Context: ctx}, i)
	if err != nil {
		return fmt.Errorf("get ip addr: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("index", i.Uint64()).
		Str("ipAddr", addr).
		Msg("sending vote to aggregator")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}

	client := NewOracleNodeClient(conn)
	_, err = client.SendVote(ctx, &SendVoteRequest{
		Index:     v.index,
		Request:   event.Request.Uint64(),
		BlockHash: hash.Bytes(),
		Signature: sig,
	})
	if err != nil {
		return fmt.Errorf("send vote: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Msg("received response")

	return nil
}
