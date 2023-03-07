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

func VerifySig(m map[string]crypto.PubKey, sig []byte, msgbytes []byte, source []byte) (bool, error) {

	//fmt.Println(m.Source)
	publickey := m[string(source)]
	//fmt.Println(source)
	if publickey == nil {
		panic("none")
	}

	return publickey.Verify(msgbytes, sig)
}
