package zkOracle

import (
	"context"
	"encoding/hex"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
)

func (n *Node) SendVote(ctx context.Context, request *SendVoteRequest) (*SendVoteResponse, error) {
	logger.Info().
		Uint64("requestNumber", request.Request).
		Str("BlockHash", hex.EncodeToString(request.BlockHash)).
		Msg("received vote")

	account, err := n.state.ReadAccount(request.Index)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "read account: %v", err)
	}

	sig := new(eddsa.Signature)
	_, err = sig.SetBytes(request.Signature)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "bytes to sig: %v", err)
	}

	err = n.votePool.add(&Vote{
		Index:     request.Index,
		Request:   big.NewInt(int64(request.Request)),
		BlockHash: common.BytesToHash(request.BlockHash),
		Sender:    *account.PublicKey,
		Signature: *sig,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add to pool: %v", err)
	}
	return &SendVoteResponse{}, nil
}

func (n *Node) mustEmbedUnimplementedOracleNodeServer() {
	panic("implement me")
}
