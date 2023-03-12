package main

import (
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/http_api"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/interface/storage"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/usecase"
)

func main() {
	// dependencies
	db := &storage.InMemoryDB{}

	// business logic
	createUC := usecase.NewCreate(db)
	closeUC := usecase.NewClose(db)
	listUC := usecase.NewList(db)
	getUC := usecase.NewGet(db)
	transferUC := usecase.NewTransfer(db)
	withdrawUC := usecase.NewWithdraw(db)
	depositUC := usecase.NewDeposit(db)

	// REST:
	controller := http_api.NewController(createUC, closeUC, listUC, getUC, transferUC, withdrawUC, depositUC)
	api := http_api.NewAPI(controller)

	if err := api.Listen("localhost:8080"); err != nil {
		panic(err)
	}

	// MessageBased: setup and register handlers
	//controller := rpc.NewMessageController(createUC, closeUC, listUC, getUC, transferUC, withdrawUC, depositUC)
	//msgHandler := rpc.NewMessageHandlerAdapter()
	//msgHandler.RegisterMessageHandler(rpc.CreateAccount, controller.HandleCreateAccount)
	//msgHandler.RegisterMessageHandler(rpc.CloseAccount, controller.HandleCloseAccount)
	//msgHandler.RegisterMessageHandler(rpc.GetAccountNumbers, controller.HandleGetAccountNumbers)
	//msgHandler.RegisterMessageHandler(rpc.GetAccount, controller.HandleGetAccount)
	//msgHandler.RegisterMessageHandler(rpc.Transfer, controller.HandleTransfer)
	//msgHandler.RegisterMessageHandler(rpc.Withdraw, controller.HandleWithdraw)
	//msgHandler.RegisterMessageHandler(rpc.Deposit, controller.HandleDeposit)
	// inject & start generic api wrapper
	//genericAPI := message_api.NewMessageBasedAPI(message_api.CreateCodecProvider(message_api.JSONCodec{}), msgHandler.Handle)
	//genericAPI.Listen("localhost:8080")
}
