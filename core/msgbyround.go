package core

import (
	"bytes"
	"errors"
	"strconv"
)

type Round struct {
	Msgs        [N][]Message
	Roundnumber int
}

func (r *Round) GetMsgByRef(hash []byte) (Message, error) {

	for _, msg := range r.Msgs {
		for _, m := range msg {
			if bytes.Equal(m.GetHash(), hash) {
				return m, nil
			}
		}
	}
	return nil, errors.New("no such message for" + string(hash) + "in round" + strconv.Itoa(r.Roundnumber))

}

func (r *Round) GetMsgByRefsBatch(hashes [][]byte) ([]Message, error) {

	msgs := make([]Message, 0)
	searchMap := make(map[string]bool)
	for _, hash := range hashes {
		searchMap[string(hash)] = true
	}

	for _, msg := range r.Msgs {
		for _, m := range msg {
			if searchMap[string(m.GetHash())] {
				msgs = append(msgs, m)
			}
		}
	}

	return msgs, nil
}
