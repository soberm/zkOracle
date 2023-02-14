package zkOracle

import (
	"context"
	"encoding/hex"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
)

func (n *Node) SendVote(ctx context.Context, request *SendVoteRequest) (*SendVoteResponse, error) {
	logger.Info().
		Uint64("requestNumber", request.Request).
		Str("blockHash", hex.EncodeToString(request.BlockHash)).
		Msg("received vote")
	err := n.votePool.add(&Vote{
		index:     request.Request,
		request:   big.NewInt(int64(request.Request)),
		blockHash: request.BlockHash,
		sender:    eddsa.PublicKey{},
		signature: eddsa.Signature{},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add to pool: %w", err)
	}
	return &SendVoteResponse{}, nil
}

func (n *Node) mustEmbedUnimplementedOracleNodeServer() {
	panic("implement me")
}
