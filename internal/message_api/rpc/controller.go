package rpc

import (
	"errors"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/message_api/model"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/usecase"
)

type MessageController struct {
	createUC   *usecase.Create
	closeUC    *usecase.Close
	listUC     *usecase.List
	getUC      *usecase.Get
	transferUC *usecase.Transfer
	withdrawUC *usecase.Withdraw
	depositUC  *usecase.Deposit
}

func NewMessageController(createUC *usecase.Create, closeUC *usecase.Close, listUC *usecase.List, getUC *usecase.Get, transferUC *usecase.Transfer, withdrawUC *usecase.Withdraw, depositUC *usecase.Deposit) *MessageController {
	return &MessageController{createUC: createUC, closeUC: closeUC, listUC: listUC, getUC: getUC, transferUC: transferUC, withdrawUC: withdrawUC, depositUC: depositUC}
}

func (c MessageController) HandleCreateAccount(v interface{}) interface{} {
	createReqModel, _ := v.(*model.CreateAccountRequestModel)

	output, err := c.createUC.Do(usecase.CreateInput{Owner: createReqModel.Owner})
	if err != nil {
		return handleError(err)
	}

	return model.CreateAccountResponseModel{ID: output.ID}
}

func (c MessageController) HandleCloseAccount(v interface{}) interface{} {
	closeReqModel, _ := v.(*model.CloseAccountRequestModel)

	output, err := c.closeUC.Do(usecase.CloseInput{ID: closeReqModel.ID})
	if err != nil {
		return handleError(err)
	}

	return model.CloseAccountResponseModel{Closed: output.Closed}
}

func (c MessageController) HandleGetAccountNumbers(v interface{}) interface{} {
	output, err := c.listUC.Do()
	if err != nil {
		return handleError(err)
	}

	return model.GetAccountNumbersResponseModel{Accounts: output.AccountsNumbers}
}

func (c MessageController) HandleGetAccount(v interface{}) interface{} {
	getReqModel, _ := v.(*model.GetAccountRequestModel)

	output, err := c.getUC.Do(usecase.GetInput{ID: getReqModel.ID})
	if err != nil {
		return handleError(err)
	}

	return model.GetAccountResponseModel{
		ID:       output.ID,
		Balance:  output.Balance,
		Owner:    output.Owner,
		IsActive: output.IsActive,
	}
}

func (c MessageController) HandleTransfer(v interface{}) interface{} {
	transferReqModel, _ := v.(*model.TransferRequestModel)

	_, err := c.transferUC.Do(usecase.TransferInput{
		FromID: transferReqModel.FromID,
		ToID:   transferReqModel.ToID,
		Amount: transferReqModel.Amount,
	})
	if err != nil {
		return handleError(err)
	}

	return model.TransferResponseModel{}
}

func (c MessageController) HandleWithdraw(v interface{}) interface{} {
	reqModel, _ := v.(*model.WithdrawRequestModel)

	err := c.withdrawUC.Do(usecase.WithdrawInput{
		ID:     reqModel.ID,
		Amount: reqModel.Amount,
	})
	if err != nil {
		return handleError(err)
	}

	return model.WithdrawResponseModel{}
}

func (c MessageController) HandleDeposit(v interface{}) interface{} {
	reqModel, _ := v.(*model.DepositRequestModel)

	err := c.depositUC.Do(usecase.DepositInput{
		ID:     reqModel.ID,
		Amount: reqModel.Amount,
	})
	if err != nil {
		return handleError(err)
	}

	return model.DepositResponseModel{}
}

func handleError(err error) error {
	target := &usecase.Error{}
	if errors.As(err, &target) {
		return &model.ErrorResponseModel{
			Code:    target.Error(),
			Message: target.Msg(),
		}
	}
	return &model.ErrorResponseModel{
		Code:    "unknown",
		Message: "unknown error occurred",
	}
}
