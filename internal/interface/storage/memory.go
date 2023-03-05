package storage

import "github.com/c0depwn/fhnw-vesys-bank-server/internal/domain"

type InMemoryDB struct {
	bank *domain.Bank
}

func (db *InMemoryDB) Get() (*domain.Bank, error) {
	if db.bank == nil {
		db.bank = &domain.Bank{Accounts: map[string]*domain.Account{}}
	}
	return db.bank, nil
}
