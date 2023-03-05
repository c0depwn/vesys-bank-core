package usecase

import (
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/domain"
	"github.com/c0depwn/fhnw-vesys-bank-server/pkg/slices"
	"golang.org/x/exp/maps"
)

type ListOutput struct {
	AccountsNumbers []string
}

type List struct {
	db DB
}

func NewList(db DB) *List {
	return &List{db: db}
}

func (uc List) Do() (ListOutput, error) {
	bank, err := uc.db.Get()
	if err != nil {
		return ListOutput{}, err
	}

	activeAccounts := slices.Filter(maps.Values(bank.Accounts), func(item *domain.Account) bool { return item.IsActive })
	ids := slices.Map(activeAccounts, func(item *domain.Account) string { return item.ID })

	return ListOutput{AccountsNumbers: ids}, nil
}
