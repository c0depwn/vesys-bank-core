package usecase

type CloseInput struct {
	ID string
}

type CloseOutput struct {
	Closed bool
}

type Close struct {
	db DB
}

func NewClose(db DB) *Close {
	return &Close{db: db}
}

func (uc Close) Do(in CloseInput) (CloseOutput, error) {
	bank, err := uc.db.Get()
	if err != nil {
		return CloseOutput{}, err
	}

	account, exists := bank.Accounts[in.ID]
	if !exists {
		return CloseOutput{Closed: false}, nil
	}

	if !account.IsActive {
		return CloseOutput{Closed: false}, nil
	}
	if account.Balance != 0 {
		return CloseOutput{Closed: false}, nil
	}

	account.Lock()
	defer account.Unlock()

	account.IsActive = false

	return CloseOutput{Closed: true}, nil
}
