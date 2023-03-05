package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/api/model"
)

type Type int

const (
	CreateAccount Type = iota
	CloseAccount
	GetAccountNumbers
	GetAccount
	Transfer
	Withdraw
	Deposit
	Error
)

func (t Type) RequestInstance() interface{} {
	switch t {
	case CreateAccount:
		return &model.CreateAccountRequestModel{}
	case CloseAccount:
		return &model.CloseAccountRequestModel{}
	case GetAccountNumbers:
		return &model.GetAccountNumbersRequestModel{}
	case GetAccount:
		return &model.GetAccountRequestModel{}
	case Transfer:
		return &model.TransferRequestModel{}
	case Withdraw:
		return &model.WithdrawRequestModel{}
	case Deposit:
		return &model.DepositRequestModel{}
	}

	panic(fmt.Errorf("unsupported message type %v", t))
}

type Message struct {
	Type Type            `json:"type"`
	Data json.RawMessage `json:"data"`
}
