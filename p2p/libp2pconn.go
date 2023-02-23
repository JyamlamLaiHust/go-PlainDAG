package p2p

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type MsgWithSig struct {
	Msg interface{}
	Sig []byte
}

type NetworkDealer struct {
	connPool     map[string]*conn
	msgch        chan MsgWithSig
	h            host.Host
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex

	timeout           time.Duration
	ctx               context.Context
	ctxCancel         context.CancelFunc
	ctxLock           sync.RWMutex
	reflectedTypesMap map[uint8]reflect.Type
}

type conn struct {
	dest   string
	w      *bufio.Writer
	encode *codec.Encoder
}

/*
write me some code to serialize the struct NetworkDealer
*/

func main() {

}

func (n *NetworkDealer) Listen() {
	listenStream := func(s network.Stream) {
		log.Println("Received a connection")
		r := bufio.NewReader(s)
		go n.HandleConn(r)
	}
	n.h.SetStreamHandler(protocol.ID("PlainDAG"), listenStream)

}

func (n *NetworkDealer) HandleConn(r *bufio.Reader) error {
	rpcType, err := r.ReadByte()
	dec := codec.NewDecoder(r, &codec.MsgpackHandle{})
	if err != nil {
		return err
	}
	reflectedType, ok := n.reflectedTypesMap[rpcType]
	if !ok {
		return errors.New(fmt.Sprintf("type of the msg (%d) is unknown", rpcType))
	}
	msgBody := reflect.Zero(reflectedType).Interface()
	if err := dec.Decode(&msgBody); err != nil {
		return err
	}

	var sig []byte
	if err := dec.Decode(&sig); err != nil {
		return err
	}

	msgWithSig := MsgWithSig{
		Msg: msgBody,
		Sig: sig,
	}
	select {
	case n.msgch <- msgWithSig:
	case <-n.shutdownCh:
		return errors.New("shut down")
	}
	return nil
}

func (n *NetworkDealer) Connect(port int, addr string) {

}

func (n *NetworkDealer) ConnectWithMultiaddr(multi string) {

}
