package core

import (
	"encoding/json"
	"flag"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/PlainDAG/go-PlainDAG/config"
	"github.com/PlainDAG/go-PlainDAG/p2p"
)

type Node struct {
	//DAG ledger structure
	bc      *Blockchain `json:"bc"`
	network *p2p.NetworkDealer
	handler Messagehandler
	//thread-safe integer
	currentround atomic.Uint32 `json:"currentround"`

	cfg *config.Config
}

func (n *Node) GenTrans() {

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

				go n.HandleMsg(msgAsserted)
			}
		}

	}
}

func (n *Node) HandleMsg(msg Message) {
	if err := n.handler.HandleMsg(msg); err != nil {
		panic(err)
	}
}

func (n *Node) ConnecttoOthers() error {
	err := n.network.Connectpeers(n.cfg.Id, n.cfg.IdaddrMap, n.cfg.IdportMap, n.cfg.Pubkeyothersmap)
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
		time.Sleep(10 * time.Millisecond)

		ref := Ref{
			Index: 1,
			H:     []byte{1, 2, 3},
		}

		refs := make([]Ref, 0)
		refs = append(refs, ref)
		msg, err := NewMroundmsg(1, refs, n.cfg.Pubkey)
		if err != nil {
			panic(err)
		}
		msgbytes, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		sig, err := n.cfg.Prvkey.Sign(msgbytes)
		if err != nil {
			panic(err)
		}
		err = n.SendMsgToAll(2, msgbytes, sig)
		if err != nil {
			panic(err)
		}
	}

}

func NewNode(filepath string) (*Node, error) {
	c := config.Loadconfig(filepath)
	n, err := p2p.Startpeer(c.Port, c.Prvkey, ReflectedTypesMap)
	if err != nil {
		return nil, err
	}

	node := Node{
		bc:      NewBlockchain(),
		network: n,
	}
	node.cfg = c
	node.handler = NewStatichandler(&node)
	node.currentround.Store(0)

	return &node, err
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
	time.Sleep(10 * time.Second)
	err = n.ConnecttoOthers()
	if err != nil {
		return nil, err
	}
	return n, nil

}
