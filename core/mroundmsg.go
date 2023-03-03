package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
)

func (m *BasicMsg) DisplayinJson() error {

	b, _ := json.Marshal(m)

	fmt.Println(string(b))
	return nil
}

func (m *BasicMsg) Encode() ([]byte, error) {
	h := sha256.Sum256([]byte(fmt.Sprintf("%v", m)))
	return h[:], nil
}

func (m *BasicMsg) GetRN() int {
	return m.Rn
}

func (m *BasicMsg) GetRefs() [][]byte {
	return m.References
}

func (m *BasicMsg) GetHash() []byte {
	return m.Hash
}

func (m *BasicMsg) GetSource() []byte {
	return m.Source
}

func (m *BasicMsg) VerifySig(n *Node, sig []byte) (bool, error) {
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

func (m *BasicMsg) HavePath(msg Message, rounds []*Round, targetround *Round) (bool, error) {
	// hashes, indexes := m.GetRefs()
	refs := m.GetRefs()
	for _, round := range rounds {
		msgs, err := round.getMsgByRefsBatch(refs)
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

		for k, v := range uniqueRefs {
			if v {

				trueRefs = append(trueRefs, []byte(k))
			}
		}
		refs = trueRefs

	}
	msgtocheck, err := targetround.getMsgByRefsBatch(refs)
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

func (m *BasicMsg) VerifyFields(n *Node) error {
	if len(m.References) < 4*f+1 {
		return errors.New("not enough references")
	}
	if n.cfg.StringpubkeyMap[string(m.Source)] == nil {
		return errors.New("no such public key")
	}
	newm := BasicMsg{
		Rn:         m.Rn,
		References: m.References,
		Source:     m.Source,
		plaintext:  messageconst,
	}
	hash, err := newm.Encode()
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, m.Hash) {
		return errors.New("hash not match")
	}
	return nil

}

func NewBasicMsg(rn int, refs [][]byte, source []byte) (*BasicMsg, error) {

	m := BasicMsg{
		Rn:         rn,
		References: refs,
		Source:     source,
		plaintext:  messageconst,
	}

	var err error
	m.Hash, err = m.Encode()
	return &m, err
}
