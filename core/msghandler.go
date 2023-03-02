package core

import "errors"

type Messagehandler interface {
	HandleMsg(msg Message, sig []byte) error
}

type Statichandler struct {
	n          *Node
	futureVers map[int]*Round
}

func (sh *Statichandler) HandleMsg(msg Message, sig []byte) error {
	b, err := msg.VerifySig(sh.n, sig)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("signature verification failed")
	}
	msg.DisplayinJson()
	return nil
}

func NewStatichandler(n *Node) *Statichandler {
	return &Statichandler{
		n:          n,
		futureVers: make(map[int]*Round),
	}
}
