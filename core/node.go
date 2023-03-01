package core

import (
	"sync/atomic"
	"time"

	"github.com/PlainDAG/go-PlainDAG/p2p"
)

type Node struct {
	//DAG ledger structure
	bc      *Blockchain `json:"bc"`
	network *p2p.NetworkDealer
	handler Messagehandler
	//thread-safe integer
	currentround atomic.Uint32 `json:"currentround"`
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
				msgAsserted.SetSource(msg.Source)
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
	err := n.network.Connectpeers()
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
		time.Sleep(5 * time.Second)

		ref := Ref{
			Index: 1,
			H:     []byte{1, 2, 3},
		}

		refs := make([]Ref, 0)
		refs = append(refs, ref)
		msg, err := NewMroundmsg(1, refs, "source")
		if err != nil {
			panic(err)
		}
		err = n.SendMsgToAll(2, msg, []byte{1, 2, 3})
		if err != nil {
			panic(err)
		}
	}

}
func NewNode(filepath string) *Node {
	n, err := p2p.Startpeer(filepath, ReflectedTypesMap)
	if err != nil {
		panic(err)
	}

	node := Node{
		bc:      NewBlokchain(),
		network: n,
	}
	node.handler = NewStatichandler(&node)
	node.currentround.Store(0)

	return &node
}
