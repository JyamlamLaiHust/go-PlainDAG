package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

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

func (m *Mroundmsg) GetRefs() [][]byte {
	return m.References
}

func (m *Mroundmsg) GetHash() []byte {
	return m.Hash
}

func (m *Mroundmsg) VerifySig(n *Node, sig []byte) (bool, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	//fmt.Println(m.Source)
	publickey := n.cfg.StringpubkeyMap[string(m.Source)]
	if publickey == nil {
		panic("none")
	}

	return publickey.Verify(bytes, sig)
}

// msg is the target message to be checked
// msgbyrounds are the messages whose round number is less than message m but larger than the target message
// targetmsground is the messageround whose round number is equal to the target message

func (m *Mroundmsg) HavePath(msg Message, rounds []*Round, targetround *Round) (bool, error) {
	// hashes, indexes := m.GetRefs()
	refs := m.GetRefs()
	for _, round := range rounds {
		msgs, err := round.GetMsgByRefsBatch(refs)
		if err != nil {
			panic(err)
		}
		uniqueRefs := make(map[string]bool)
		for _, m := range msgs {
			refs := m.GetRefs()
			for _, ref := range refs {
				uniqueRefs[string(ref)] = true
			}
		}

		trueRefs := make([][]byte, 0)
		// convert string to byte array

		for k, v := range uniqueRefs {
			if v {

				trueRefs = append(trueRefs, []byte(k))
			}
		}
		refs = trueRefs

	}
	msgtocheck, err := targetround.GetMsgByRefsBatch(refs)
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

func NewMroundmsg(rn uint32, refs [][]byte, source []byte) (*Mroundmsg, error) {
	m := Mroundmsg{
		Rn:         rn,
		References: refs,
		Source:     source,
	}

	var err error
	m.Hash, err = m.Encode()
	return &m, err
}
