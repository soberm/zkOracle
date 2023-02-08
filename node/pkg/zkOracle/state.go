package zkOracle

import (
	"fmt"
	"hash"
)

type State struct {
	hFunc hash.Hash
	data  []byte
	hData []byte
}

func NewState(hFunc hash.Hash, accounts []*Account) (*State, error) {

	data := make([]byte, accountSize*nbAccounts)
	hData := make([]byte, hFunc.Size()*nbAccounts)

	for i, account := range accounts {

		hFunc.Reset()

		accountData := account.Serialize()
		_, err := hFunc.Write(accountData)
		if err != nil {
			return nil, fmt.Errorf("hash account: %w", err)
		}
		s := hFunc.Sum(nil)

		copy(data[i*accountSize:(i+1)*accountSize], accountData)
		copy(hData[i*hFunc.Size():(i+1)*hFunc.Size()], s)
	}

	return &State{
		hFunc: hFunc,
		data:  data,
		hData: hData,
	}, nil
}

func (s *State) UpdateState(vote *Vote) error {
	account, err := s.ReadAccount(vote.index)
	if err != nil {
		return fmt.Errorf("read account: %w", err)
	}
	fmt.Printf("%v", account)
	return nil
}

func (s *State) ReadAccount(i uint64) (Account, error) {
	var res Account
	res.Deserialize(s.data[int(i)*accountSize : int(i)*accountSize+accountSize])
	return res, nil
}

func (s *State) Data() []byte {
	return s.data
}

func (s *State) HashData() []byte {
	return s.hData
}
