package core

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Messagehandler interface {
	handleMsg(msg Message, sig []byte) error
	getFutureMsgByRound(rn int) []Message
	tryHandle(msg Message) error
	handleFutureVers(rn int) error

	buildContextForRound(rn int)
	signalFutureVersHandled(rn int)
	readyForRound(rn int)
}

type Statichandler struct {
	n              *Node
	futureVers     map[int][]Message
	futureVerslock sync.RWMutex

	waitingChanMap     map[int]chan bool
	waitingChanMaplock sync.RWMutex

	readyToSendMap     map[int]chan bool
	readyToSendMapLock sync.RWMutex

	isDoneWithFutureVers map[int]chan bool
	isDoneWithFuturelock sync.RWMutex

	isSent     map[int]bool
	isSentLock sync.RWMutex
}

func (sh *Statichandler) signalFutureVersHandled(rn int) {
	sh.isDoneWithFuturelock.Lock()

	ch := sh.isDoneWithFutureVers[rn]
	ch <- true
	sh.isDoneWithFuturelock.Unlock()

}
func (sh *Statichandler) readyForRound(rn int) {
	if rn == 1 {
		return
	}

	// sh.waitingChanMaplock.Lock()
	// chwaiting := sh.waitingChanMap[rn-1]
	// sh.waitingChanMaplock.Unlock()
	//fmt.Println("are you stuck here in waiting?")

	// var wg sync.WaitGroup
	// for i := 0; i < 4*f; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		<-chwaiting
	// 		fmt.Println("done")
	// 		wg.Done()
	// 	}()
	// }

	// wg.Wait()
	sh.readyToSendMapLock.Lock()
	chready := sh.readyToSendMap[rn-1]
	sh.readyToSendMapLock.Unlock()

	<-chready

	//fmt.Println("are you stuck here in waiting?")
	//fmt.Println("done waiting for round " + strconv.Itoa(rn-1))
	sh.isDoneWithFuturelock.RLock()
	chfuturedone := sh.isDoneWithFutureVers[rn-1]

	sh.isDoneWithFuturelock.RUnlock()
	//fmt.Println("are you stuck here in ready for round?")
	<-chfuturedone
	//fmt.Println("are you stuck here in ready for round?")

	close(chfuturedone)

}
func (sh *Statichandler) buildContextForRound(rn int) {
	sh.addWaitingChan(rn)
	sh.addIsDoneChan(rn)
	sh.addReadyToSendChan(rn)

}

func (sh *Statichandler) addIsDoneChan(rn int) {
	sh.isDoneWithFuturelock.Lock()
	sh.isDoneWithFutureVers[rn] = make(chan bool, 1)
	sh.isDoneWithFuturelock.Unlock()
}

func (sh *Statichandler) addWaitingChan(rn int) {
	sh.waitingChanMaplock.Lock()
	sh.waitingChanMap[rn] = make(chan bool, 5*f)
	sh.waitingChanMaplock.Unlock()
}

func (sh *Statichandler) addReadyToSendChan(rn int) {
	sh.readyToSendMapLock.Lock()
	sh.readyToSendMap[rn] = make(chan bool, 1)
	sh.readyToSendMapLock.Unlock()
}

func (sh *Statichandler) handleMsg(msg Message, sig []byte) error {
	// err := sh.VerifyandCheckMsg(msg, sig)
	// if err != nil {
	// 	return err
	// }

	isFuture := sh.storeFutureVers(msg)
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
	//fmt.Println("ends here tryhandle1?")
	sh.waitingChanMaplock.Lock()
	ch := sh.waitingChanMap[rn]

	sh.waitingChanMaplock.Unlock()
	//fmt.Println("ends here tryhandle1?")

	//fmt.Println("ends here tryhandle2?")
	sh.readyToSendMapLock.Lock()
	if len(ch) < 4*f {
		ch <- true
	} else {
		//fmt.Println("now 4f")
		chready := sh.readyToSendMap[rn]

		sh.isSentLock.Lock()

		if !sh.isSent[rn] {
			//fmt.Println("here?")
			chready <- true

			sh.isSent[rn] = true
			sh.isSentLock.Unlock()
			close(chready)
			close(ch)
			sh.readyToSendMapLock.Unlock()
			fmt.Println("handle msg success from    " + strconv.Itoa(id) + "round number: " + strconv.Itoa(rn))
			return nil
		}
		sh.isSentLock.Unlock()
		//chready <- true

	}
	fmt.Println("handle msg success from    " + strconv.Itoa(id) + "round number: " + strconv.Itoa(rn))

	sh.readyToSendMapLock.Unlock()
	return nil
}

func (sh *Statichandler) storeFutureVers(msg Message) bool {
	sh.futureVerslock.Lock()
	//fmt.Println("stuck here?")
	if msg.GetRN() > int(sh.n.currentround.Load()) {
		//sh.futureVerslock.Lock()
		sh.futureVers[msg.GetRN()] = append(sh.futureVers[msg.GetRN()], msg)
		sh.futureVerslock.Unlock()
		return true
	}
	sh.futureVerslock.Unlock()
	return false
}

func (sh *Statichandler) handleFutureVers(rn int) error {

	msgsNextRound := sh.getFutureMsgByRound(rn)
	if msgsNextRound == nil {
		sh.signalFutureVersHandled(rn)
		//fmt.Println("signaled")
		return nil
	}
	//fmt.Println(len(msgsNextRound))
	var err error
	var wg sync.WaitGroup
	//fmt.Println("are you stuck here?")

	for _, msg := range msgsNextRound {
		m := msg
		_, canattach := sh.n.bc.GetRound(rn - 1).checkCanAttach(msg)
		if canattach {
			wg.Add(1)
			//fmt.Println("wg add by 1")
			go func() {
				//fmt.Println("hi")
				err = sh.tryHandle(m)
				//fmt.Println("hiagain")
				//fmt.Println("wg decrease by 1")
				wg.Done()
			}()
		} else {
			go sh.tryHandle(msg)
		}
	}
	wg.Wait()
	//fmt.Println("are you stuck here?")
	sh.signalFutureVersHandled(rn)

	return err
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

func (sh *Statichandler) getFutureMsgByRound(rn int) []Message {
	sh.futureVerslock.RLock()
	defer sh.futureVerslock.RUnlock()
	msgs := sh.futureVers[rn]
	return msgs
}
func NewStatichandler(n *Node) *Statichandler {
	return &Statichandler{
		n:                    n,
		futureVers:           make(map[int][]Message),
		waitingChanMap:       make(map[int]chan bool),
		isDoneWithFutureVers: make(map[int]chan bool),
		readyToSendMap:       make(map[int]chan bool),
		isSent:               make(map[int]bool),
	}
}
