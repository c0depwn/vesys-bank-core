package main

import (
	"fmt"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/api"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/api/rpc"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/interface/storage"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/usecase"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", "0.0.0.0:8080")

	if err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
	defer listener.Close()

	log.Printf("[INFO] server is listening on %v", listener.Addr())

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

	// setup and register handlers
	controller := rpc.NewMessageController(createUC, closeUC, listUC, getUC, transferUC, withdrawUC, depositUC)
	msgHandler := rpc.NewMessageHandlerAdapter()
	msgHandler.RegisterMessageHandler(rpc.CreateAccount, controller.HandleCreateAccount)
	msgHandler.RegisterMessageHandler(rpc.CloseAccount, controller.HandleCloseAccount)
	msgHandler.RegisterMessageHandler(rpc.GetAccountNumbers, controller.HandleGetAccountNumbers)
	msgHandler.RegisterMessageHandler(rpc.GetAccount, controller.HandleGetAccount)
	msgHandler.RegisterMessageHandler(rpc.Transfer, controller.HandleTransfer)
	msgHandler.RegisterMessageHandler(rpc.Withdraw, controller.HandleWithdraw)
	msgHandler.RegisterMessageHandler(rpc.Deposit, controller.HandleDeposit)

	// inject & start generic api wrapper
	genericAPI := api.NewMessageBasedAPI(api.CreateCodecProvider(api.JSONCodec{}), msgHandler.Handle)
	genericAPI.Listen(listener)
}
