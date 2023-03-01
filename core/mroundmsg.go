package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func (m *Mroundmsg) MarshalJSON() ([]byte, error) {

	hexstrings := make([]string, len(m.ReferencesHash))
	for i, v := range m.ReferencesHash {
		hexstrings[i] = fmt.Sprintf("%x", v)
	}
	hexstring := fmt.Sprintf("%x", m.Hash)
	return json.Marshal(struct {
		RN              uint32
		ReferencesHash  []string
		ReferencesIndex []uint8
		Source          string
		Hash            string
	}{
		RN:              m.RN,
		ReferencesHash:  hexstrings,
		ReferencesIndex: m.ReferencesIndex,
		Source:          m.Source,
		Hash:            hexstring,
	})

}

func (m *Mroundmsg) DisplayinJson() error {
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
	return nil
}

func (m *Mroundmsg) Encode() ([]byte, error) {
	h := sha256.Sum256([]byte(fmt.Sprintf("%v", m)))
	return h[:], nil
}

func (m *Mroundmsg) GetRefs() ([][]byte, []uint8) {
	return m.ReferencesHash, m.ReferencesIndex
}

func (m *Mroundmsg) GetRN() uint32 {
	return m.RN
}
func (m *Mroundmsg) IsEqual(msg Message) (bool, error) {
	m2, ok := msg.(*Mroundmsg)
	if !ok {
		return false, fmt.Errorf("wrong type")
	}
	if m.RN != m2.RN {
		return false, nil
	}
	if len(m.ReferencesHash) != len(m2.ReferencesHash) {
		return false, nil
	}
	if len(m.ReferencesIndex) != len(m2.ReferencesIndex) {
		return false, nil
	}
	for i, v := range m.ReferencesHash {
		if bytes.Equal(v, m2.ReferencesHash[i]) {
			return false, nil
		}
	}
	for i, v := range m.ReferencesIndex {
		if v != m2.ReferencesIndex[i] {
			return false, nil
		}
	}
	if m.Source != m2.Source {
		return false, nil
	}
	return true, nil
}

func (m *Mroundmsg) HavePath(msg Message, msgbyrounds []*MSGByRound) (bool, error) {
	_, indexes := m.GetRefs()

	for _, msgbyround := range msgbyrounds {
		msgs := msgbyround.GetMsgByindexes(indexes)
		for _, msg := range msgs {
	}
	return false, nil

	// 	return true

}

func NewMroundmsg(rn uint32, hashset [][]byte, indexset []uint8, src string) (*Mroundmsg, error) {
	m := Mroundmsg{
		RN:              rn,
		ReferencesHash:  hashset,
		ReferencesIndex: indexset,
		Source:          src,
	}
	var err error
	m.Hash, err = m.Encode()
	return &m, err

}
