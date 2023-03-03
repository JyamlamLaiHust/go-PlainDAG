package core

import (
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
	currentround atomic.Uint32 `json:"currentround"`

	cfg *config.Config

	isSent     map[int]bool
	isSentLock sync.Mutex
}

func (n *Node) GenTrans(rn int) {
	n.isSentLock.Lock()
	if n.isSent[rn] {
		n.isSentLock.Unlock()
		return
	}
	n.isSent[rn] = true
	n.isSentLock.Unlock()
	log.Println("generate transaction for round" + strconv.Itoa(rn))
	//generate transaction

}

func (n *Node) GenFroundMsg(rn int) (*Froundmsg, error) {
	return nil, nil
}

func (n *Node) GenLroundMsg(rn int) (*Lroundmsg, error) {
	return nil, nil
}

func (n *Node) GenBasicMsg(rn int) (*BasicMsg, error) {
	return nil, nil
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

				go n.HandleMsg(msgAsserted, msg.Sig, msg.Source)
			}
		}

	}
}

func (n *Node) HandleMsg(msg Message, sig []byte, source crypto.PubKey) {
	if err := n.handler.HandleMsg(msg, sig); err != nil {
		panic(err)
	}
}

func (n *Node) ConnecttoOthers() error {
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
		time.Sleep(1000 * time.Millisecond)

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
	err = n.ConnecttoOthers()
	if err != nil {
		return nil, err
	}
	n.constructpubkeyMap()
	// get the pubkey of my own host

	return n, nil

}
