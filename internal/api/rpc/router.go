package rpc

import (
	"fmt"
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/api"
)

type SpecificMessageHandler func(v interface{}) interface{}

// TODO: middleware types
//type GenericMessageHandler func(msg Message) (Message, error)
//type MessageMiddleware func(next GenericMessageHandler) GenericMessageHandler

type MessageHandlerAdapter struct {
	handlers map[Type]SpecificMessageHandler
	//middlewares []MessageMiddleware
}

func NewMessageHandlerAdapter() MessageHandlerAdapter {
	return MessageHandlerAdapter{handlers: map[Type]SpecificMessageHandler{}}
}

// TODO: register middleware
//func (h MessageHandlerAdapter) RegisterMiddleware(f MessageMiddleware) {
//	h.middlewares = append(h.middlewares, f)
//}

func (h MessageHandlerAdapter) RegisterMessageHandler(t Type, f SpecificMessageHandler) {
	h.handlers[t] = f
}

// TODO: init middleware
//func (h MessageHandlerAdapter) Initialize() {
//	var next GenericMessageHandler
//
//	for i := range h.middlewares {
//		next = h.middlewares[i](next)
//	}
//}

func (h MessageHandlerAdapter) Handle(codec api.Codec, data []byte) ([]byte, error) {
	var rpcMsg Message
	if err := codec.Decode(data, &rpcMsg); err != nil {
		return nil, err
	}

	v := rpcMsg.Type.RequestInstance()
	if err := codec.Decode(rpcMsg.Data, &v); err != nil {
		return nil, err
	}

	handler, exists := h.handlers[rpcMsg.Type]
	if !exists {
		return nil, fmt.Errorf("no handler available for the specified message type")
	}

	rpcResponseMsg := handler(v)
	responseType := rpcMsg.Type
	if _, ok := rpcResponseMsg.(error); ok {
		responseType = Error
	}

	rawResponseMessageData, err := codec.Encode(&rpcResponseMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode rpc message data: %v", err)
	}

	rawResponseMessage, err := codec.Encode(&Message{
		Type: responseType,
		Data: rawResponseMessageData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode rpc message: %v", err)
	}

	return rawResponseMessage, nil
}
