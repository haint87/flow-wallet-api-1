package postgres

import (
	"github.com/eqlabs/flow-nft-wallet-service/pkg/store"
	"github.com/google/uuid"
)

type DataStore struct {
	store.AccountStore
}

type AccountStore struct{}

func NewDataStore() (*DataStore, error) {
	return &DataStore{
		AccountStore: NewAccountStore(),
	}, nil
}

func NewAccountStore() *AccountStore {
	return &AccountStore{}
}

func (s *AccountStore) Account(id uuid.UUID) (store.Account, error) {
	panic("not implemented") // TODO: Implement
}

func (s *AccountStore) Accounts() ([]store.Account, error) {
	panic("not implemented") // TODO: Implement
}

func (s *AccountStore) CreateAccount(a *store.Account) error {
	panic("not implemented") // TODO: Implement
}

func (s *AccountStore) DeleteAccount(id uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}