package usecase

import "fmt"

type WithdrawInput struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type Withdraw struct {
	db DB
}

func NewWithdraw(db DB) *Withdraw {
	return &Withdraw{db: db}
}

func (uc Withdraw) Do(input WithdrawInput) error {
	bank, err := uc.db.Get()
	if err != nil {
		return err
	}

	account, exists := bank.Accounts[input.ID]
	if !exists {
		return fmt.Errorf("account not found")
	}

	if !account.IsActive {
		return fmt.Errorf("account inactive")
	}

	if input.Amount < 0 {
		return &ErrInvalid
	}

	if account.Balance < input.Amount {
		return &ErrOverdraw
	}

	account.Lock()
	defer account.Unlock()

	account.Balance -= input.Amount

	return nil
}
