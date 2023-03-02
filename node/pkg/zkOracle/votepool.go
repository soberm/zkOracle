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
	results    map[uint64]map[common.Hash]uint64
	sink       chan uint64
}

func NewVotePool() *VotePool {
	return &VotePool{
		votes:      make(map[common.Hash]*Vote),
		voteHashes: make(map[uint64][]*common.Hash),
		results:    make(map[uint64]map[common.Hash]uint64),
		sink:       make(chan uint64),
	}
}

func (vp *VotePool) add(vote *Vote) error {
	vp.Lock()
	defer vp.Unlock()

	logger.Info().
		Uint64("index", vote.Index).
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

	requestID := vote.Request.Uint64()

	vp.voteHashes[requestID] = append(vp.voteHashes[vote.Request.Uint64()], &voteHash)
	vp.votes[voteHash] = vote

	_, ok = vp.results[requestID]
	if !ok {
		vp.results[requestID] = make(map[common.Hash]uint64)
	}

	_, ok = vp.results[requestID][vote.BlockHash]
	if !ok {
		vp.results[requestID][vote.BlockHash] = 0
	}
	vp.results[requestID][vote.BlockHash] += 1

	if vp.results[requestID][vote.BlockHash] == NumAccounts {
		select {
		case vp.sink <- vote.Request.Uint64():
		default:
		}
	}

	return nil
}

func (vp *VotePool) getVotes(requestID uint64) ([]*Vote, error) {
	vp.RLock()
	defer vp.RUnlock()

	votes := make([]*Vote, 0)
	for _, vote := range vp.voteHashes[requestID] {
		votes = append(votes, vp.votes[*vote])
	}

	return votes, nil
}

func (vp *VotePool) verifyVote(vote *Vote) (bool, error) {
	logger.Info().
		Uint64("requestNumber", vote.Request.Uint64()).
		Str("blockHash", vote.BlockHash.String()).
		Uint64("index", vote.Index).
		Msg("verify vote")

	hFunc := mimc.NewMiMC()
	hFunc.Write(vote.Serialize())

	isValid, err := vote.Sender.Verify(vote.Signature.Bytes(), hFunc.Sum(nil), hFunc)
	if err != nil {
		return isValid, fmt.Errorf("verify Signature: %w", err)
	}
	return isValid, nil
}
