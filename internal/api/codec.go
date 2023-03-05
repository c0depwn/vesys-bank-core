package api

import (
	"encoding/json"
	"fmt"
)

type Codec interface {
	Capability() MessageDataEncoding
	Decode(data []byte, v interface{}) error
	Encode(v interface{}) ([]byte, error)
}

type CodecProvider func(encoding MessageDataEncoding) (Codec, error)

func CreateCodecProvider(codecs ...Codec) CodecProvider {
	providerMap := map[MessageDataEncoding]Codec{}
	for i := range codecs {
		providerMap[codecs[i].Capability()] = codecs[i]
	}
	return func(encoding MessageDataEncoding) (Codec, error) {
		codec, exists := providerMap[encoding]
		if !exists {
			return nil, fmt.Errorf("unsupported encoding %v", encoding)
		}
		return codec, nil
	}
}

type JSONCodec struct{}

func (J JSONCodec) Capability() MessageDataEncoding {
	return JSON
}

func (J JSONCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (J JSONCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
