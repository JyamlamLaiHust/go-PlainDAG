package core

import (
	"errors"
	"sync"
)

type Messagehandler interface {
	HandleMsg(msg Message, sig []byte) error
}

type Statichandler struct {
	n              *Node
	futureVers     map[int]*Round
	futureVerslock sync.RWMutex
}

func (sh *Statichandler) HandleMsg(msg Message, sig []byte) error {
	err := sh.VerifyandCheckMsg(msg, sig)
	if err != nil {
		return err
	}
	//msg.DisplayinJson()
	return nil
}

func (sh *Statichandler) VerifyandCheckMsg(msg Message, sig []byte) error {
	b, err := msg.VerifySig(sh.n, sig)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("signature verification failed")
	}

	if err := msg.VerifyFields(sh.n); err != nil {
		return err
	}
	return nil
}

func NewStatichandler(n *Node) *Statichandler {
	return &Statichandler{
		n:          n,
		futureVers: make(map[int]*Round),
	}
}
