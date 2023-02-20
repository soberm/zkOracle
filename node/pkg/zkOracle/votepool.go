package zkOracle

import (
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
)

type VotePool struct {
	sync.RWMutex
	votes      map[common.Hash]*Vote
	voteHashes map[uint64][]*common.Hash
	sink       chan uint64
}

func NewVotePool() *VotePool {
	return &VotePool{
		votes:      make(map[common.Hash]*Vote),
		voteHashes: make(map[uint64][]*common.Hash),
		sink:       make(chan uint64),
	}
}

func (vp *VotePool) add(vote *Vote) error {
	vp.Lock()
	defer vp.Unlock()

	logger.Info().
		Uint64("Index", vote.Index).
		Uint64("requestNumber", vote.Request.Uint64()).
		Msg("adding vote")

	isValid, err := vp.verifyVote(vote)
	if err != nil {
		return fmt.Errorf("verify vote: %w", err)
	}

	if !isValid {
		return errors.New("invalid vote")
	}

	voteHash := crypto.Keccak256Hash(vote.Serialize())
	_, ok := vp.votes[voteHash]
	if ok {
		return fmt.Errorf("vote already exists")
	}

	vp.voteHashes[vote.Request.Uint64()] = append(vp.voteHashes[vote.Request.Uint64()], &voteHash)
	vp.votes[voteHash] = vote

	if len(vp.voteHashes[vote.Request.Uint64()]) == nbAccounts {
		select {
		case vp.sink <- vote.Request.Uint64():
		default:
		}
	}

	return nil
}

func (vp *VotePool) getVotes(requestID uint64) ([]*Vote, error) {

	votes := make([]*Vote, 0)
	for _, vote := range vp.voteHashes[requestID] {
		votes = append(votes, vp.votes[*vote])
	}

	return votes, nil
}

func (vp *VotePool) verifyVote(vote *Vote) (bool, error) {
	logger.Info().
		Uint64("requestNumber", vote.Request.Uint64()).
		Str("BlockHash", vote.BlockHash.String()).
		Uint64("Index", vote.Index).
		Msg("verify vote")

	hasher := mimc.NewMiMC()
	hasher.Write(vote.Serialize())
	msg := hasher.Sum(nil)

	isValid, err := vote.Sender.Verify(vote.Signature.Bytes(), msg, mimc.NewMiMC())
	if err != nil {
		return isValid, fmt.Errorf("verify Signature: %w", err)
	}
	return isValid, nil
}
