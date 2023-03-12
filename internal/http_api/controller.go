package http_api

import (
	"encoding/json"
	"errors"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/message_api/model"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/usecase"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	createUC   *usecase.Create
	closeUC    *usecase.Close
	listUC     *usecase.List
	getUC      *usecase.Get
	transferUC *usecase.Transfer
	withdrawUC *usecase.Withdraw
	depositUC  *usecase.Deposit
}

func NewController(createUC *usecase.Create, closeUC *usecase.Close, listUC *usecase.List, getUC *usecase.Get, transferUC *usecase.Transfer, withdrawUC *usecase.Withdraw, depositUC *usecase.Deposit) *Controller {
	return &Controller{createUC: createUC, closeUC: closeUC, listUC: listUC, getUC: getUC, transferUC: transferUC, withdrawUC: withdrawUC, depositUC: depositUC}
}

func (c *Controller) HandleCreateAccount(res http.ResponseWriter, req *http.Request) {
	var createReqModel model.CreateAccountRequestModel
	if err := readJSON(req.Body, &createReqModel); err != nil {
		writeError(res, http.StatusBadRequest, err)
		return
	}

	output, err := c.createUC.Do(usecase.CreateInput{Owner: createReqModel.Owner})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusCreated, &model.CreateAccountResponseModel{ID: output.ID})
}

func (c *Controller) HandleCloseAccount(res http.ResponseWriter, req *http.Request) {
	accountID := mux.Vars(req)["id"]
	output, err := c.closeUC.Do(usecase.CloseInput{ID: accountID})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.CloseAccountResponseModel{Closed: output.Closed})
}

func (c *Controller) HandleGetAccountNumbers(res http.ResponseWriter, req *http.Request) {
	output, err := c.listUC.Do()
	if err != nil {
		writeError(res, http.StatusBadRequest, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.GetAccountNumbersResponseModel{Accounts: output.AccountsNumbers})
}

func (c *Controller) HandleGetAccount(res http.ResponseWriter, req *http.Request) {
	accountID := mux.Vars(req)["id"]
	output, err := c.getUC.Do(usecase.GetInput{ID: accountID})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.GetAccountResponseModel{
		ID:       output.ID,
		Balance:  output.Balance,
		Owner:    output.Owner,
		IsActive: output.IsActive,
	})
}

func (c *Controller) HandleTransfer(res http.ResponseWriter, req *http.Request) {
	var transferReqModel model.TransferRequestModel
	if err := readJSON(req.Body, &transferReqModel); err != nil {
		writeError(res, http.StatusBadRequest, err)
		return
	}

	_, err := c.transferUC.Do(usecase.TransferInput{
		FromID: transferReqModel.FromID,
		ToID:   transferReqModel.ToID,
		Amount: transferReqModel.Amount,
	})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.TransferResponseModel{})
}

func (c *Controller) HandleWithdraw(res http.ResponseWriter, req *http.Request) {
	accountID := mux.Vars(req)["id"]

	var reqModel model.WithdrawRequestModel
	if err := readJSON(req.Body, &reqModel); err != nil {
		writeError(res, http.StatusBadRequest, err)
		return
	}

	err := c.withdrawUC.Do(usecase.WithdrawInput{
		ID:     accountID,
		Amount: reqModel.Amount,
	})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.WithdrawResponseModel{})
}

func (c *Controller) HandleDeposit(res http.ResponseWriter, req *http.Request) {
	accountID := mux.Vars(req)["id"]

	var reqModel model.DepositRequestModel
	if err := readJSON(req.Body, &reqModel); err != nil {
		writeError(res, http.StatusBadRequest, err)
		return
	}

	err := c.depositUC.Do(usecase.DepositInput{
		ID:     accountID,
		Amount: reqModel.Amount,
	})
	if err != nil {
		writeError(res, http.StatusNotAcceptable, err)
		return
	}

	writeJSON(res, http.StatusOK, &model.DepositResponseModel{})
}

func writeError(writer http.ResponseWriter, status int, err error) {
	log.Printf("[ERROR] status=%v err=%v", status, err)

	target := &usecase.Error{}
	if errors.As(err, &target) {
		writeJSON(writer, status, &model.ErrorResponseModel{
			Code:    target.Error(),
			Message: target.Msg(),
		})
		return
	}
	writeJSON(writer, status, &model.ErrorResponseModel{
		Code:    "unknown",
		Message: "unknown error occurred",
	})
}

func readJSON(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

func writeJSON(writer http.ResponseWriter, status int, v interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(v); err != nil {
		log.Printf("[Error] writeJSON: %v", err)
		panic(err)
	}
}
