package message_api

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
)

type MessageDataEncoding = uint16

const (
	JSON = MessageDataEncoding(iota)
)

type Metadata struct {
	id       uuid.UUID
	size     uint16
	encoding uint16
	status   uint16
}

func (meta Metadata) GetEncoding() MessageDataEncoding {
	return meta.encoding
}

type Message struct {
	Meta Metadata
	Body []byte
}

type MessageReader func(reader io.Reader) (Message, error)
type MessageWriter func(writer io.Writer, msg Message) error
type MessageHandler func(codec Codec, data []byte) ([]byte, error)

type MessageBased struct {
	provider  CodecProvider
	msgReader MessageReader
	msgWriter MessageWriter
	handler   MessageHandler
}

func NewMessageBasedAPI(
	provider CodecProvider,
	handler MessageHandler,
) MessageBased {
	return MessageBased{
		provider:  provider,
		msgReader: readMessage,
		msgWriter: writeMessage,
		handler:   handler,
	}
}

func (api *MessageBased) Listen(addr string) {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
	defer listener.Close()

	log.Printf("[INFO] server is listening on %v", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] failed to accept new connection: %v", err)
			continue
		}

		log.Printf("[INFO] accepted new connection from %v", conn.RemoteAddr())
		go api.handle(conn)
	}
}

func (api *MessageBased) handle(rw io.ReadWriteCloser) {
	for {
		log.Printf("[INFO] waiting for message...")

		receivedMsg, err := api.msgReader(rw)
		if err != nil {
			log.Printf("[ERROR] failed to read message: %v", err)
			rw.Close()
			return
		}

		codec, err := api.provider(receivedMsg.Meta.GetEncoding())
		if err != nil {
			log.Printf("[ERROR] invalid message encoding: %v", err)
			rw.Close()
			return
		}

		log.Printf("[INFO] handling message %v...", receivedMsg.Meta.id)
		responseData, err := api.handler(codec, receivedMsg.Body)
		if err != nil {
			log.Printf("[ERROR] failed to process message: %v", err)
			rw.Close()
			return
		}
		log.Printf("[INFO] handled message: %v", string(responseData))

		responseMsg := Message{
			Meta: Metadata{
				id:       receivedMsg.Meta.id,
				size:     uint16(len(responseData)),
				encoding: receivedMsg.Meta.encoding,
				status:   0,
			},
			Body: responseData,
		}

		if err := api.msgWriter(rw, responseMsg); err != nil {
			log.Printf("[ERROR] failed to write message: %v", err)
			rw.Close()
			return
		}

		log.Printf("[INFO] request handled successfully")
	}
}
