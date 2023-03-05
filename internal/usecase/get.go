package usecase

import "fmt"

type GetInput struct {
	ID string
}

type GetOutput struct {
	ID       string
	Balance  float64
	Owner    string
	IsActive bool
}

type Get struct {
	db DB
}

func NewGet(db DB) *Get {
	return &Get{db: db}
}

func (uc Get) Do(input GetInput) (GetOutput, error) {
	bank, err := uc.db.Get()
	if err != nil {
		return GetOutput{}, err
	}

	account, exists := bank.Accounts[input.ID]
	if !exists {
		return GetOutput{}, fmt.Errorf("account %v not found", input.ID)
	}

	return GetOutput{
		ID:       account.ID,
		Balance:  account.Balance,
		Owner:    account.Owner,
		IsActive: account.IsActive,
	}, nil
}
