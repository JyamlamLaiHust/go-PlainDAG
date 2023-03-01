package core

import (
	"bytes"
	"errors"
	"strconv"
)

type MSGByRound struct {
	Msgs        [N][]Message
	Roundnumber int
}

func (mr *MSGByRound) GetMsgByRef(ref *Ref) (Message, error) {
	if mr.Msgs[ref.Index] == nil {
		return nil, errors.New("no message at index" + strconv.Itoa(int(ref.Index)))
	}

	for _, msg := range mr.Msgs[ref.Index] {
		if bytes.Equal(msg.GetHash(), ref.H) {
			return msg, nil
		}
	}
	return nil, errors.New("no such message for" + string(ref.H) + "at index" +
		strconv.Itoa(int(ref.Index)) + "in round" + strconv.Itoa(mr.Roundnumber))

}

func (mr *MSGByRound) GetMsgByRefsBatch(refs []Ref) ([]Message, error) {
	msgs := make([]Message, len(refs))
	var err error
	for i, r := range refs {
		msgs[i], err = mr.GetMsgByRef(&r)
		if err != nil {
			return nil, err
		}
	}
	return msgs, nil
}
