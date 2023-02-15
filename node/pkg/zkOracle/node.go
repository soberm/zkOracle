package zkOracle

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"math/big"
	"math/rand"
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
	server          *grpc.Server
}

func NewNode() (*Node, error) {

	ethClient, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		return nil, fmt.Errorf("dial eth: %w", err)
	}

	contract, err := NewZKOracleContract(common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"), ethClient)
	if err != nil {
		return nil, fmt.Errorf("oracle contract: %w", err)
	}

	r := rand.New(rand.NewSource(0))
	eddsaPrivateKey, err := eddsa.GenerateKey(r)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	validator := NewValidator(ethClient, contract, eddsaPrivateKey)
	aggregator := NewAggregator()

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		return nil, fmt.Errorf("ecdsa private key: %w", err)
	}

	return &Node{
		server:          grpc.NewServer(),
		contract:        contract,
		validator:       validator,
		aggregator:      aggregator,
		chainID:         chainID,
		ecdsaPrivateKey: ecdsaPrivateKey,
		eddsaPrivateKey: eddsaPrivateKey,
		ethClient:       ethClient,
	}, nil
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
