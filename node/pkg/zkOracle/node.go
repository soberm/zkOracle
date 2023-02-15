package zkOracle

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"math/big"
	"net"
)

type Node struct {
	UnimplementedOracleNodeServer
	chainID         *big.Int
	ethClient       *ethclient.Client
	ecdsaPrivateKey *ecdsa.PrivateKey
	eddsaPrivateKey *eddsa.PrivateKey
	contract        *ZKOracleContract
	aggregator      *Aggregator
	validator       *Validator
	votePool        *VotePool
	state           *State
	stateSync       *StateSync
	server          *grpc.Server
}

func NewNode(config *Config) (*Node, error) {

	ethClient, err := ethclient.Dial(config.Ethereum.TargetAddress)
	if err != nil {
		return nil, fmt.Errorf("dial eth: %w", err)
	}

	contract, err := NewZKOracleContract(common.HexToAddress(config.ContractAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("oracle contract: %w", err)
	}

	b, err := hex.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("eddsa private key to bytes: %w", err)
	}

	eddsaPrivateKey := new(eddsa.PrivateKey)
	_, err = eddsaPrivateKey.SetBytes(b)
	if err != nil {
		return nil, fmt.Errorf("eddsa private key from bytes: %w", err)
	}

	validator := NewValidator(config.Index, ethClient, contract, eddsaPrivateKey)
	aggregator := NewAggregator()

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(config.Ethereum.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("ecdsa private key: %w", err)
	}

	state, err := NewState(mimc.NewMiMC(), nil)
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	stateSync := NewStateSync(state, contract)

	server := grpc.NewServer()

	node := &Node{
		server:          server,
		contract:        contract,
		validator:       validator,
		aggregator:      aggregator,
		chainID:         chainID,
		ecdsaPrivateKey: ecdsaPrivateKey,
		eddsaPrivateKey: eddsaPrivateKey,
		ethClient:       ethClient,
		state:           state,
		stateSync:       stateSync,
		votePool:        NewVotePool(),
	}
	RegisterOracleNodeServer(server, node)

	return node, nil
}

func (n *Node) Register(ipAddr string) error {
	auth, err := bind.NewKeyedTransactorWithChainID(n.ecdsaPrivateKey, n.chainID)
	if err != nil {
		return fmt.Errorf("new transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(20000000000)

	x := new(big.Int)
	y := new(big.Int)

	n.eddsaPrivateKey.PublicKey.A.X.ToBigIntRegular(x)
	n.eddsaPrivateKey.PublicKey.A.Y.ToBigIntRegular(y)

	logger.Info().
		Str("pubKeyX", x.String()).
		Str("pubKeyY", y.String()).
		Str("ipAddr", ipAddr).
		Msg("register")

	tx, err := n.contract.Register(auth, ZKOraclePublicKey{
		X: x,
		Y: y,
	}, ipAddr)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	_, err = bind.WaitMined(context.Background(), n.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	return nil
}

func (n *Node) Run(listener net.Listener) error {

	go func() {
		if err := n.stateSync.Synchronize(); err != nil {
			logger.Err(err).Msg("synchronize")
		}
	}()

	err := n.Register(listener.Addr().String())
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	go func() {
		err := n.aggregator.Aggregate(context.Background())
		if err != nil {
			logger.Err(err).Msg("aggregate")
		}
	}()

	go func() {
		err := n.validator.Validate(context.Background())
		if err != nil {
			logger.Err(err).Msg("validate")
		}
	}()

	return n.server.Serve(listener)
}

func (n *Node) Stop() {
	n.server.Stop()
}