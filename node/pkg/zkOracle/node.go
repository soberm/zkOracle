package zkOracle

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type Node struct {
	UnimplementedOracleNodeServer
	aggregator *Aggregator
	validator  *Validator
	votePool   *VotePool
	server     *grpc.Server
}

func NewNode() (*Node, error) {
	return &Node{server: grpc.NewServer()}, nil
}

func (n *Node) Run(listener net.Listener) error {

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
