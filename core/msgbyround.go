package core

type MSGByRound struct {
	Msgs        [N]Message
	Roundnumber int
}

func (mr *MSGByRound) GetMsgByindex(index uint8) Message {
	return mr.Msgs[index]
}

func (mr *MSGByRound) GetMsgByindexes(indexes []uint8) []Message {
	var msgs []Message
	for _, index := range indexes {
		msgs = append(msgs, mr.Msgs[index])
	}
	return msgs
}
