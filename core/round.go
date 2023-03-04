package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Round struct {
	msgs        [5*f + 1][]Message
	roundNumber int
	messageLock sync.RWMutex

	//checkmap is used to check whether a message has been received
	//1. upon a message is received, we first check if these's a bool channel for this messgage.hash in the checkmap
	// if not, it means there's no messages in higher round waiting for this message
	// if so, we put a a bool value into the channel
	// 2.  when we receive a message, we first check if all the messages it references have been stored in the DAG
	// if true, we simply store it in the DAG.
	// if false, we find the hashes that are not in the DAG, and use these hashes to construct the map and let the channel block this message
	// until all messages it references are received
	checkMap     map[string]chan bool
	checkMapLock sync.RWMutex
	isFround     bool
}

func (r *Round) getMsgByRef(hash []byte) (Message, error) {
	r.messageLock.RLock()
	defer r.messageLock.RUnlock()
	for _, msg := range r.msgs {
		for _, m := range msg {
			if bytes.Equal(m.GetHash(), hash) {
				return m, nil
			}
		}
	}
	return nil, errors.New("no such message for" + string(hash) + "in round" + strconv.Itoa(r.roundNumber))

}

func (r *Round) getMsgByRefsBatch(hashes [][]byte) ([]Message, error) {
	r.messageLock.RLock()
	defer r.messageLock.RUnlock()
	msgs := make([]Message, 0)
	searchMap := make(map[string]bool)
	for _, hash := range hashes {
		searchMap[string(hash)] = true
	}

	for _, msg := range r.msgs {
		for _, m := range msg {
			if searchMap[string(m.GetHash())] {
				msgs = append(msgs, m)
			}
		}
	}

	return msgs, nil
}

func (r *Round) checkCanAttach(m Message) ([][]byte, bool) {
	r.messageLock.RLock()

	refs := m.GetRefs()
	searchMap := make(map[string]bool)
	for _, ref := range refs {
		searchMap[string(ref)] = true
	}

	for _, msg := range r.msgs {
		for _, m := range msg {
			if searchMap[string(m.GetHash())] {
				delete(searchMap, string(m.GetHash()))
			}
		}
	}
	r.messageLock.RUnlock()
	if len(searchMap) == 0 {
		return nil, true
	}
	missingRefs := make([][]byte, 0)
	for k, _ := range searchMap {
		missingRefs = append(missingRefs, []byte(k))
	}

	return missingRefs, false
}

func (r *Round) tryAttach(m Message, currentRound *Round, id int) {
	// CheckCanAttack checks whether all the msgs message m references are stored in the last round

	r.checkMapLock.Lock()
	missingrefs, canattach := r.checkCanAttach(m)
	// if true, attack msg m to the current round
	if canattach {
		currentRound.attachMsg(m, id)
		r.checkMapLock.Unlock()
		return
	}
	// r.checkMapLock.Lock()
	for _, ref := range missingrefs {
		if _, ok := r.checkMap[string(ref)]; !ok {
			r.checkMap[string(ref)] = make(chan bool)
		}
	}
	r.checkMapLock.Unlock()
	var wg sync.WaitGroup

	for _, ref := range missingrefs {
		wg.Add(1)
		rcopy := make([]byte, len(ref))
		copy(rcopy, ref)
		c := r.checkMap[string(rcopy)]
		go func() {
			b := <-c
			c <- b
			wg.Done()
		}()

	}
	wg.Wait()
	currentRound.checkMapLock.Lock()

	if _, ok := currentRound.checkMap[string(m.GetHash())]; !ok {
		//currentRound.checkMap[string(m.GetHash())] = make(chan bool)
		currentRound.attachMsg(m, id)
		return
	}
	currentRound.checkMap[string(m.GetHash())] <- true
	currentRound.checkMapLock.Unlock()
	go func() {
		time.Sleep(2 * time.Second)
		currentRound.checkMapLock.Lock()

		close(currentRound.checkMap[string(m.GetHash())])
		currentRound.checkMapLock.Unlock()
	}()
	currentRound.attachMsg(m, id)

}

func (r *Round) attachMsg(m Message, id int) {
	r.messageLock.Lock()

	//fmt.Println(id)
	r.msgs[id] = append(r.msgs[id], m)
	r.messageLock.Unlock()

}

func (r *Round) retMsgsToRef() [][]byte {
	msgsByte := make([][]byte, 0)
	r.messageLock.RLock()
	for _, msg := range r.msgs {

		if msg == nil {
			continue
		}
		fmt.Println(len(msg))
		for _, m := range msg {

			if m == nil {
				continue
			}
			msgsByte = append(msgsByte, m.GetHash())
		}

	}
	r.messageLock.RUnlock()
	return msgsByte
}
func newRound(rn int, m Message, id int) (*Round, error) {
	if m.GetRN() != rn {
		return nil, errors.New("round number not match")
	}
	var msglists [5*f + 1][]Message
	for i := 0; i < 5*f+1; i++ {
		msglists[i] = make([]Message, 0)
	}
	msglists[id] = append(msglists[id], m)
	return &Round{
		msgs:         msglists,
		roundNumber:  rn,
		checkMap:     make(map[string]chan bool),
		messageLock:  sync.RWMutex{},
		checkMapLock: sync.RWMutex{},
	}, nil
}

// this function returns the first round(round 0) in the DAG and is hardcoded
func newStaticRound() (*Round, error) {
	var msglists [5*f + 1][]Message
	for i := 0; i < 5*f+1; i++ {
		msglists[i] = make([]Message, 0)
		if msg, err := NewBasicMsg(0, [][]byte{}, []byte("Source"+strconv.Itoa(i))); err != nil {
			return nil, err
		} else {
			msglists[i] = append(msglists[i], msg)
		}

	}
	return &Round{
		msgs:         msglists,
		roundNumber:  0,
		checkMap:     make(map[string]chan bool),
		messageLock:  sync.RWMutex{},
		checkMapLock: sync.RWMutex{},
	}, nil

}
