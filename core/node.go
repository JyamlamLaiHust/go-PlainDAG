package core

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/JyamlamLaiHUST/go-PlainDAG/config"
	"github.com/JyamlamLaiHUST/go-PlainDAG/p2p"
	"github.com/JyamlamLaiHUST/go-PlainDAG/sign"
	"github.com/JyamlamLaiHUST/go-PlainDAG/utils"
)

type Node struct {
	//DAG ledger structure
	bc      *Blockchain `json:"bc"`
	network *p2p.NetworkDealer
	handler Messagehandler
	//thread-safe integer
	currentround atomic.Int64 `json:"currentround"`

	cfg *config.Config

	ls        *LeaderSelector
	committer *StaticCommitter
}

func (n *Node) genTrans(rn int) (Message, error) {
	if rn < 3 {
		return n.genBasicMsg(rn)
	}
	if rn%100 == 0 {
		log.Println("generate transaction for round" + strconv.Itoa(rn))
	} //generate transaction
	if rn%rPerwave == 0 {

		return n.genBasicMsg(rn)
	} else if rn%rPerwave == 1 {
		return n.genBasicMsg(rn)
	} else if rn%rPerwave == 2 {
		return n.genLroundMsg(rn)
	}
	return nil, nil
}

func (n *Node) genFroundMsg(rn int) (*FRoundMsg, error) {
	return nil, nil
}

func (n *Node) genLroundMsg(rn int) (*LRoundMsg, error) {
	basic, err := n.genBasicMsg(rn)
	if err != nil {
		return nil, err
	}

	lround := n.bc.GetRound(rn - 2)
	mround := n.bc.GetRound(rn - 1)
	bytes, err := lround.genArefs(basic, mround)
	if err != nil {
		return nil, err
	}

	lroundmsg, err := NewLroundMsg(bytes, basic)

	if err != nil {
		return nil, err
	}
	indexes, err := lround.getIndexByRefsBatch(bytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("Generated lround message at round ", rn, "and A-references ", indexes)
	return lroundmsg, nil

}

func (n *Node) genBasicMsg(rn int) (*BasicMsg, error) {
	lastRound := n.bc.GetRound(rn - 1)
	if lastRound == nil {
		return nil, errors.New("last round is nil")
	}
	//generate transaction
	refsByte := lastRound.retMsgsToRef()
	//fmt.Println(len(msgsByte))
	basicMsg, err := NewBasicMsg(rn, refsByte, n.cfg.Pubkeyraw)
	if err != nil {
		return nil, err
	}
	//fmt.Println("ends here?")
	return basicMsg, nil

}

func (n *Node) genThresMsg(rn int) *ThresSigMsg {

	bytes := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(bytes, int64(rn))
	//fmt.Println("generated for round  ", rn, "    ", bytes)
	s := sign.SignTSPartial(n.cfg.TSPrvKey, bytes)
	thresSigMsg := &ThresSigMsg{
		Wn:     rn / rPerwave,
		Sig:    s,
		Source: n.cfg.Pubkeyraw,
	}
	//fmt.Println(thresSigMsg.source)
	return thresSigMsg

}
func (n *Node) paceToNextRound() (Message, error) {
	//generate transaction
	rn := int(n.currentround.Load())
	n.handler.buildContextForRound(rn + 1)

	//this removal is only used to save memory when the code is not finished.
	// if rn > 11 {
	// 	minustenRound := n.bc.GetRound(rn - 10)
	// 	minustenRound.rmvAllMsgsWhenCommitted()
	// }
	msg, err := n.genTrans(rn + 1)
	if err != nil {
		return nil, err
	}

	msgbytes, sig, err := utils.MarshalAndSign(msg, n.cfg.Prvkey)
	if err != nil {
		return nil, err
	}
	if rn < 3 {

		go n.SendMsgToAll(1, msgbytes, sig)
	} else {
		msgtype := (rn + 1) % 3
		//fmt.Println(msgtype)
		//fmt.Println((rn + 1) % 3)
		go n.SendMsgToAll(uint8(msgtype), msgbytes, sig)
	}
	if rn%rPerwave == 1 && rn != 1 {
		thresSigMsg := n.genThresMsg(rn + 1)

		//fmt.Println(thresSigMsg.sig, thresSigMsg.source, thresSigMsg.wn)
		thresSigMsgBytes, sig, err := utils.MarshalAndSign(thresSigMsg, n.cfg.Prvkey)
		n.handler.handleThresMsg(thresSigMsg, sig, thresSigMsgBytes)
		if err != nil {
			return nil, err
		}
		//rintln(thresSigMsgBytes)
		go n.SendMsgToAll(3, thresSigMsgBytes, sig)
	}

	//initialize a new round with the newly generated message msg
	newR, err := newRound(rn+1, msg, n.cfg.Id)

	if err != nil {
		return nil, err
	}
	n.bc.AddRound(newR)
	msg.AfterAttach(n)
	n.currentround.Add(1)

	go n.handler.handleFutureVers(rn + 1)
	//n.SendMsgToAll(1, msgbytes, sig)
	return msg, err
}

func (n *Node) HandleMsgForever() {
	for {
		select {

		case msg := <-n.network.ExtractMsg():
			//log.Println("receive msg: ", msg.Msg)
			switch msgAsserted := msg.Msg.(type) {
			case Message:
				go n.handleMsg(msgAsserted, msg.Sig, msg.Msgbytes)
			case *ThresSigMsg:
				// fmt.Println("received thresmsg")
				//fmt.Println(msg.Msg)
				go n.handleThresMsg(msgAsserted, msg.Sig, msg.Msgbytes)
			}

		}

	}
}

func (n *Node) handleMsg(msg Message, sig []byte, msgbytes []byte) {
	if err := n.handler.handleMsg(msg, sig, msgbytes); err != nil {
		panic(err)
	}
}

func (n *Node) handleThresMsg(msg *ThresSigMsg, sig []byte, msgbytes []byte) {
	if err := n.handler.handleThresMsg(msg, sig, msgbytes); err != nil {
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
	if err := n.network.Broadcast(messagetype, msg, sig, n.cfg.Simlatency); err != nil {
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

func (n *Node) serialize(filepath string) error {
	//Todo
	//serialize the committed messages to the database
	return nil
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
