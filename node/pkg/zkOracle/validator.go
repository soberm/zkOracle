package zkOracle

import (
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const CONFIRMATIONS uint64 = 5

type Validator struct {
	ethClient        *ethclient.Client
	zkOracleContract *ZKOracleContract
	privateKey       *eddsa.PrivateKey
}

func NewValidator(ethClient *ethclient.Client, zkOracleContract *ZKOracleContract, privateKey *eddsa.PrivateKey) *Validator {
	return &Validator{ethClient: ethClient, zkOracleContract: zkOracleContract, privateKey: privateKey}
}

func (v *Validator) Validate() error {
	return nil
}

func (v *Validator) WatchAndHandleBlockRequestedEvent(ctx context.Context) error {
	sink := make(chan *ZKOracleContractBlockRequested)
	defer close(sink)

	sub, err := v.zkOracleContract.WatchBlockRequested(&bind.WatchOpts{
		Context: context.Background(),
	}, sink, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			if err := v.HandleBlockRequestedEvent(ctx, event); err != nil {
				fmt.Printf("handle BlockRequestedEvent: %v", err)
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (v *Validator) HandleBlockRequestedEvent(ctx context.Context, event *ZKOracleContractBlockRequested) error {
	block, err := v.ethClient.HeaderByNumber(ctx, event.Number)
	if err != nil {
		return fmt.Errorf("block by number: %w", err)
	}

	currentBlockNumber, err := v.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("blocknumber: %w", err)
	}

	if currentBlockNumber-block.Number.Uint64() < CONFIRMATIONS {
		return fmt.Errorf("block not confirmed")
	}

	sig, err := v.privateKey.Sign(block.Hash().Bytes(), mimc.NewMiMC())
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}
	fmt.Printf("%v", sig)

	i, err := v.zkOracleContract.GetAggregator(&bind.CallOpts{
		Context: ctx},
	)
	if err != nil {
		return fmt.Errorf("get aggregator: %w", err)
	}

	addr, err := v.zkOracleContract.GetIPAddress(&bind.CallOpts{Context: ctx}, i)
	if err != nil {
		return fmt.Errorf("get ip addr: %w", err)
	}
	fmt.Printf("%v", addr)
	//TODO: Send message

	return nil
}
