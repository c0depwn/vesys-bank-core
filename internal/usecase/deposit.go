package usecase

import "fmt"

type DepositInput struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type Deposit struct {
	db DB
}

func NewDeposit(db DB) *Deposit {
	return &Deposit{db: db}
}

func (uc Deposit) Do(input DepositInput) error {
	bank, err := uc.db.Get()
	if err != nil {
		return err
	}

	account, exists := bank.Accounts[input.ID]
	if !exists {
		return fmt.Errorf("account not found")
	}

	if !account.IsActive {
		return &ErrInactive
	}

	if input.Amount < 0 {
		return &ErrInvalid
	}

	account.Lock()
	defer account.Unlock()

	account.Balance += input.Amount

	return nil
}
