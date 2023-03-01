package core

type Messagehandler interface {
	HandleMsg(msg Message) error
}

type Statichandler struct {
	n          *Node
	futureVers map[int]*MSGByRound
}

func (sh *Statichandler) HandleMsg(msg Message) error {

	msg.DisplayinJson()
	return nil
}

func NewStatichandler(n *Node) *Statichandler {
	return &Statichandler{
		n:          n,
		futureVers: make(map[int]*MSGByRound),
	}
}
