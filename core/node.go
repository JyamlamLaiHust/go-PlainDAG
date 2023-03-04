package core

import (
	"errors"
	"flag"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlainDAG/go-PlainDAG/config"
	"github.com/PlainDAG/go-PlainDAG/p2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

type Node struct {
	//DAG ledger structure
	bc      *Blockchain `json:"bc"`
	network *p2p.NetworkDealer
	handler Messagehandler
	//thread-safe integer
	currentround atomic.Int64 `json:"currentround"`

	cfg *config.Config

	isSentMap  map[int]bool
	isSentLock sync.Mutex
}

func (n *Node) genTrans(rn int) (Message, error) {

	if rn%100 == 0 {
		log.Println("generate transaction for round" + strconv.Itoa(rn))
	} //generate transaction
	return n.genBasicMsg(rn)
}

func (n *Node) genFroundMsg(rn int) (*Froundmsg, error) {
	return nil, nil
}

func (n *Node) genLroundMsg(rn int) (*Lroundmsg, error) {
	return nil, nil
}

func (n *Node) genBasicMsg(rn int) (*BasicMsg, error) {
	lastRound := n.bc.GetRound(rn - 1)
	if lastRound == nil {
		return nil, errors.New("last round is nil")
	}
	//generate transaction
	msgsByte := lastRound.retMsgsToRef()
	//fmt.Println(len(msgsByte))
	basicMsg, err := NewBasicMsg(rn, msgsByte, n.cfg.Pubkeyraw)
	if err != nil {
		return nil, err
	}
	//fmt.Println("ends here?")
	return basicMsg, nil

}

func (n *Node) paceToNextRound() (Message, error) {
	//generate transaction
	rn := int(n.currentround.Load())
	n.handler.buildContextForRound(rn + 1)
	msg, err := n.genTrans(rn + 1)

	if err != nil {
		return nil, err
	}

	newR, err := newRound(rn+1, msg, n.cfg.Id)

	if err != nil {
		return nil, err
	}
	n.bc.AddRound(newR)

	n.currentround.Add(1)
	go n.handler.handleFutureVers(int(n.currentround.Load()))
	go n.SendMsgToAll(1, msg, []byte{1})
	// msgbytes, err := json.Marshal(msg)
	// if err != nil {
	// 	panic(err)
	// }

	// sig, err := n.cfg.Prvkey.Sign(msgbytes)
	// if err != nil {
	// 	panic(err)
	// }

	return msg, err
}

func (n *Node) WaitAtRn(rn int) {

}

func (n *Node) HandleMsgForever() {
	for {
		select {
		case <-n.network.ExtractShutdown():
			return
		case msg := <-n.network.ExtractMsg():
			//log.Println("receive msg: ", msg.Msg)
			switch msgAsserted := msg.Msg.(type) {
			case Message:

				go n.handleMsg(msgAsserted, msg.Sig, msg.Source)
			}
		}

	}
}

func (n *Node) handleMsg(msg Message, sig []byte, source crypto.PubKey) {
	if err := n.handler.handleMsg(msg, sig); err != nil {
		panic(err)
	}
}

func (n *Node) connecttoOthers() error {
	err := n.network.Connectpeers(n.cfg.Id, n.cfg.IdaddrMap, n.cfg.IdportMap, n.cfg.IdPubkeymap)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) SendMsg(messagetype uint8, msg interface{}, sig []byte, dest string) error {
	if err := n.network.SendMsg(messagetype, msg, sig, dest); err != nil {
		return err
	}
	return nil
}

func (n *Node) SendMsgToAll(messagetype uint8, msg interface{}, sig []byte) error {
	if err := n.network.Broadcast(messagetype, msg, sig); err != nil {
		return err
	}
	return nil
}

func (n *Node) SendForever() {

	for {
		n.handler.readyForRound(int(n.currentround.Load()) + 1)
		_, err := n.paceToNextRound()
		if err != nil {
			panic(err)
		}
		//msg.DisplayinJson()

		//time.Sleep(100 * time.Millisecond)
		// H := []byte{1, 2, 3}

		// refs := make([][]byte, 0)
		// refs = append(refs, H)

		// msg, err := NewMroundmsg(1, refs, n.cfg.Pubkeyraw)
		// if err != nil {
		// 	panic(err)
		// }
		// // for _, peer := range n.network.H.Peerstore().Peers() {
		// // 	s := peer.Pretty()

		// // 	fmt.Println(s)
		// // }
		// msgbytes, err := json.Marshal(msg)
		// if err != nil {
		// 	panic(err)
		// }

		// sig, err := n.cfg.Prvkey.Sign(msgbytes)
		// if err != nil {
		// 	panic(err)
		// }
		// err = n.SendMsgToAll(2, msg, sig)
		// if err != nil {
		// 	panic(err)
		// }
	}

}

func StartandConnect() (*Node, error) {
	index := flag.Int("f", 0, "config file path")
	flag.Parse()
	//convert int to string
	filepath := "node" + strconv.Itoa(*index)
	n, err := NewNode(filepath)
	if err != nil {
		return nil, err
	}

	time.Sleep(15 * time.Second)
	err = n.connecttoOthers()
	if err != nil {
		return nil, err
	}
	n.constructpubkeyMap()
	// get the pubkey of my own host

	return n, nil

}
