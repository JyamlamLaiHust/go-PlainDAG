package utils

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func MarshalAndSign(msg interface{}, prvkey crypto.PrivKey) ([]byte, []byte, error) {
	msgbytes, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, err
	}
	sig, err := prvkey.Sign(msgbytes)
	if err != nil {
		return nil, nil, err
	}

	return msgbytes, sig, nil
}
