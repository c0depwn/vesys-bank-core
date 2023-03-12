package model

type CreateAccountRequestModel struct {
	Owner string `json:"owner"`
}

type CreateAccountResponseModel struct {
	ID string `json:"id"`
}

type CloseAccountRequestModel struct {
	ID string `json:"id"`
}

type CloseAccountResponseModel struct {
	Closed bool `json:"closed"`
}

type GetAccountNumbersRequestModel struct {
}

type GetAccountNumbersResponseModel struct {
	Accounts []string `json:"accounts"`
}

type GetAccountRequestModel struct {
	ID string `json:"id"`
}

type GetAccountResponseModel struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Owner    string  `json:"owner"`
	IsActive bool    `json:"is_active"`
}

type TransferRequestModel struct {
	FromID string  `json:"from_id"`
	ToID   string  `json:"to_id"`
	Amount float64 `json:"amount"`
}

type TransferResponseModel struct{}

type WithdrawRequestModel struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type WithdrawResponseModel struct{}

type DepositRequestModel struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type DepositResponseModel struct{}
