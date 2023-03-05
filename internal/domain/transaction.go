package domain

import "fmt"

var (
	InvalidAmountErr = fmt.Errorf("transaction amount cannot be negative")
	OverdrawErr      = fmt.Errorf("insufficient funds")
	InactiveErr      = fmt.Errorf("account is inactive")
)

func Transfer(from *Account, to *Account, amount float64) error {
	if amount < 0 {
		return InvalidAmountErr
	}
	if from.Balance < amount {
		return OverdrawErr
	}
	if !from.IsActive || !to.IsActive {
		return InactiveErr
	}

	from.mu.Lock()
	to.mu.Lock()
	defer from.mu.Unlock()
	defer to.mu.Unlock()

	from.Balance -= amount
	to.Balance += amount
	return nil
}
