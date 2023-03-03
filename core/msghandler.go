package core

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Messagehandler interface {
	handleMsg(msg Message, sig []byte) error
	getMsgByRound(rn int) []Message
	tryHandle(msg Message) error
}

type Statichandler struct {
	n              *Node
	futureVers     map[int][]Message
	futureVerslock sync.RWMutex
}

func (sh *Statichandler) handleMsg(msg Message, sig []byte) error {
	err := sh.VerifyandCheckMsg(msg, sig)
	if err != nil {
		return err
	}

	isFuture := sh.handleFutureVers(msg)
	if isFuture {
		return nil
	}

	//msg.DisplayinJson()

	return sh.tryHandle(msg)
}

func (sh *Statichandler) tryHandle(msg Message) error {
	id := sh.n.cfg.StringIdMap[string(msg.GetSource())]

	rn := msg.GetRN()
	lastRound := sh.n.bc.GetRound(rn - 1)
	lastRound.tryAttach(msg, sh.n.bc.GetRound(rn), id)
	fmt.Println("handle msg success from    " + strconv.Itoa(id) + "round number: " + strconv.Itoa(rn))
	return nil
}
func (sh *Statichandler) handleFutureVers(msg Message) bool {
	sh.futureVerslock.Lock()
	if msg.GetRN() > int(sh.n.currentround.Load()) {
		//sh.futureVerslock.Lock()
		sh.futureVers[msg.GetRN()] = append(sh.futureVers[msg.GetRN()], msg)
		sh.futureVerslock.Unlock()
		return true
	}

	return false
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

func (sh *Statichandler) getMsgByRound(rn int) []Message {
	sh.futureVerslock.Lock()
	defer sh.futureVerslock.Unlock()
	msgs := sh.futureVers[rn]
	return msgs
}
func NewStatichandler(n *Node) *Statichandler {
	return &Statichandler{
		n:          n,
		futureVers: make(map[int][]Message),
	}
}
