package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// func (m *Mroundmsg) MarshalJSON() ([]byte, error) {

// 	return nil, nil
// }

func (m *Mroundmsg) DisplayinJson() error {
	b, _ := json.Marshal(m)

	fmt.Println(string(b))
	return nil
}

func (m *Mroundmsg) Encode() ([]byte, error) {
	h := sha256.Sum256([]byte(fmt.Sprintf("%v", m)))
	return h[:], nil
}

func (m *Mroundmsg) GetRN() uint32 {
	return m.Rn
}

func (m *Mroundmsg) GetRefs() []Ref {
	return m.References
}

func (m *Mroundmsg) GetHash() []byte {
	return m.Hash
}

// msg is the target message to be checked
// msgbyrounds are the messages whose round number is less than message m but larger than the target message
// targetmsground is the messageround whose round number is equal to the target message

func (m *Mroundmsg) HavePath(msg Message, msgbyrounds []*MSGByRound, targetmsground *MSGByRound) (bool, error) {
	// hashes, indexes := m.GetRefs()
	refs := m.GetRefs()
	for _, msgbyround := range msgbyrounds {
		msgs, err := msgbyround.GetMsgByRefsBatch(refs)
		if err != nil {
			panic(err)
		}
		uniqueRefs := make(map[*Ref]bool)
		for _, m := range msgs {
			refs := m.GetRefs()
			for _, ref := range refs {
				uniqueRefs[&ref] = true
			}
		}
		trueRefs := []Ref{}
		for k, v := range uniqueRefs {
			if v {
				trueRefs = append(trueRefs, *k)
			}
		}
		refs = trueRefs

	}
	msgtocheck, err := targetmsground.GetMsgByRefsBatch(refs)
	if err != nil {
		panic(err)
	}
	for _, m := range msgtocheck {
		if bytes.Equal(m.GetHash(), msg.GetHash()) {
			return true, nil
		}
	}
	return false, nil

}

func (m *Mroundmsg) SetSource(dest string) {
	m.Source = dest
}
func NewMroundmsg(rn uint32, refs []Ref, source string) (*Mroundmsg, error) {
	m := Mroundmsg{
		Rn:         rn,
		References: refs,
	}
	var err error
	m.Hash, err = m.Encode()
	return &m, err
}
