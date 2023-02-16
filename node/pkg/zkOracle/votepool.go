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
		Uint64("index", vote.index).
		Uint64("requestNumber", vote.request.Uint64()).
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
	
	vp.voteHashes[vote.request.Uint64()] = append(vp.voteHashes[vote.request.Uint64()], &voteHash)
	vp.votes[voteHash] = vote

	if len(vp.voteHashes[vote.request.Uint64()]) == nbAccounts {
		select {
		case vp.sink <- vote.request.Uint64():
		default:
		}
	}

	return nil
}

func (vp *VotePool) verifyVote(vote *Vote) (bool, error) {
	logger.Info().
		Uint64("requestNumber", vote.request.Uint64()).
		Str("blockHash", vote.blockHash.String()).
		Uint64("index", vote.index).
		Msg("verify vote")

	isValid, err := vote.sender.Verify(vote.signature.Bytes(), vote.Serialize(), mimc.NewMiMC())
	if err != nil {
		return isValid, fmt.Errorf("verify signature: %w", err)
	}
	return isValid, nil
}
