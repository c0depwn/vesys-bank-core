package message_api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"io"
)

/*

Message Packet: [ ID ] [ DATA SIZE ] [ ENCODING ] [ STATUS ] [ DATA ]

Message Metadata
id       uuid 	(16 bytes)
size     uint16 (2 bytes)
encoding uint16 (2 bytes)
status   uint16 (2 bytes)

Message Data
data []byte (x bytes <= 2^16-1)
*/

func readMessage(reader io.Reader) (Message, error) {
	id, err := readUUIDField(reader)
	if err != nil {
		return Message{}, err
	}

	msgDataSize, err := readUint16Field(reader)
	if err != nil {
		return Message{}, err
	}

	encoding, err := readUint16Field(reader)
	if err != nil {
		return Message{}, err
	}

	status, err := readUint16Field(reader)
	if err != nil {
		return Message{}, err
	}

	data, err := readData(reader, msgDataSize)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Meta: Metadata{
			id:       id,
			size:     msgDataSize,
			encoding: encoding,
			status:   status,
		},
		Body: data,
	}, nil
}

func readUint16Field(reader io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read uint16 field: %v", err)
	}

	return binary.BigEndian.Uint16(buf), nil
}

func readUUIDField(reader io.Reader) (uuid.UUID, error) {
	buf := make([]byte, 16)

	n, err := io.ReadFull(reader, buf)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to read UUID field: %v, n=%v", err, n)
	}

	id, err := uuid.FromBytes(buf)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID bytes: %v", err)
	}

	return id, nil
}

func readData(reader io.Reader, size uint16) ([]byte, error) {
	dataBuf := make([]byte, size)

	_, err := io.ReadFull(reader, dataBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read data field: %v", err)
	}

	return dataBuf, nil
}

func writeMessage(writer io.Writer, msg Message) error {
	msgBuf := bytes.Buffer{}

	uuidBuf, err := msg.Meta.id.MarshalBinary()
	if err != nil {
	}
	msgBuf.Write(uuidBuf)

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, msg.Meta.size)
	msgBuf.Write(buf)

	binary.BigEndian.PutUint16(buf, msg.Meta.encoding)
	msgBuf.Write(buf)

	binary.BigEndian.PutUint16(buf, msg.Meta.status)
	msgBuf.Write(buf)

	msgBuf.Write(msg.Body)

	_, err = msgBuf.WriteTo(writer)
	if err != nil {
		return fmt.Errorf("failed to write message to writer: %v", err)
	}

	return nil
}
