package usecase

import "github.com/c0depwn/fhnw-vesys-bank-server/internal/domain"

type CreateInput struct {
	Owner string
}

type CreateOutput struct {
	ID string
}

type Create struct {
	db DB
}

func NewCreate(db DB) *Create {
	return &Create{db: db}
}

func (uc Create) Do(input CreateInput) (CreateOutput, error) {
	bank, err := uc.db.Get()
	if err != nil {
		return CreateOutput{}, err
	}

	account := domain.NewAccount(input.Owner)
	bank.Accounts[account.ID] = account

	return CreateOutput{ID: account.ID}, nil
}
