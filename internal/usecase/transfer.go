package usecase

import (
	"fmt"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/domain"
)

type TransferInput struct {
	FromID string
	ToID   string
	Amount float64
}

type TransferOutput struct{}

type Transfer struct {
	db DB
}

func NewTransfer(db DB) *Transfer {
	return &Transfer{db: db}
}

func (uc Transfer) Do(input TransferInput) (TransferOutput, error) {
	bank, err := uc.db.Get()
	if err != nil {
		return TransferOutput{}, err
	}

	from, ok := bank.Accounts[input.FromID]
	if !ok {
		return TransferOutput{}, fmt.Errorf("account with id %v not found", input.FromID)
	}
	to, ok := bank.Accounts[input.ToID]
	if !ok {
		return TransferOutput{}, fmt.Errorf("account with id %v not found", input.ToID)
	}

	return TransferOutput{}, domain.Transfer(from, to, input.Amount)
}

func (uc Transfer) transfer(from *domain.Account, to *domain.Account, amount float64) error {
	if amount < 0 {
		return &ErrInvalid
	}
	if from.Balance < amount {
		return &ErrOverdraw
	}
	if !from.IsActive || !to.IsActive {
		return &ErrInactive
	}

	from.Lock()
	to.Lock()
	defer func() {
		from.Unlock()
		to.Unlock()
	}()

	from.Balance -= amount
	to.Balance += amount
	return nil
}
