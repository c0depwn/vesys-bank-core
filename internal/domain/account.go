package domain

import (
	"github.com/google/uuid"
	"sync"
)

type Account struct {
	mu sync.Mutex

	ID       string
	Balance  float64
	Owner    string
	IsActive bool
}

func NewAccount(owner string) *Account {
	return &Account{
		mu:       sync.Mutex{},
		ID:       uuid.New().String(),
		Balance:  0,
		Owner:    owner,
		IsActive: true,
	}
}

func (a *Account) Lock() {
	a.mu.Lock()
}

func (a *Account) Unlock() {
	a.mu.Unlock()
}
