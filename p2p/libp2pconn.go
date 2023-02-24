package p2p

import (
	"bufio"
	"context"

	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/libp2p/go-libp2p"

	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/PlainDAG/go-PlainDAG/config"
	"github.com/multiformats/go-multiaddr"
)

type MsgWithSig struct {
	Msg interface{}
	Sig []byte
}

type NetworkDealer struct {
	connPool     map[string]*conn
	msgch        chan MsgWithSig
	H            host.Host
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex

	// ctx               context.Context
	// ctxCancel         context.CancelFunc
	// ctxLock           sync.RWMutex

	reflectedTypesMap map[uint8]reflect.Type
	config            *config.P2pconfig
}

type conn struct {
	dest   string
	w      *bufio.Writer
	encode *codec.Encoder
}

/*
write me some code to serialize the struct NetworkDealer
*/

func makeHost(port int, prvKey crypto.PrivKey) host.Host {
	// Make the host that will handle the network requests
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	host, err := libp2p.New(libp2p.Identity(prvKey), libp2p.ListenAddrs(sourceMultiAddr))
	if err != nil {
		log.Fatal(err)
	}
	return host

}

func (n *NetworkDealer) Listen() {
	listenStream := func(s network.Stream) {
		log.Println("Received a connection")

		log.Println(s.Conn().RemotePeer().String())
		r := bufio.NewReader(s)
		go n.HandleConn(r)
	}
	n.H.SetStreamHandler(protocol.ID("PlainDAG"), listenStream)

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

func (n *NetworkDealer) Connect(port int, addr string, pubKey string) (*bufio.Writer, error) {
	return n.ConnectWithMultiaddr(PackMultiaddr(port, addr, pubKey))
}

func (n *NetworkDealer) ConnectWithMultiaddr(multi string) (*bufio.Writer, error) {

	log.Println("Connecting to ", multi)
	maddr, err := multiaddr.NewMultiaddr(multi)
	if err != nil {
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}
	n.H.Peerstore().AddAddr(info.ID, maddr, peerstore.PermanentAddrTTL)
	s, err := n.H.NewStream(context.Background(), info.ID, protocol.ID("PlainDAG"))
	if err != nil {
		return nil, err
	}
	return bufio.NewWriter(s), nil
}

func PackMultiaddr(port int, addr string, pubKey string) string {
	return fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", addr, port, pubKey)
}

func (n *NetworkDealer) SendMsg(rpcType uint8, msg interface{}, sig []byte, dest string) error {

	n.shutdownLock.Lock()
	if n.shutdown {
		n.shutdownLock.Unlock()
		return errors.New("shut down")
	}
	n.shutdownLock.Unlock()

	c, ok := n.connPool[dest]
	if !ok {
		w, err := n.ConnectWithMultiaddr(dest)
		if err != nil {
			return err
		}

		c = &conn{
			dest:   dest,
			w:      w,
			encode: codec.NewEncoder(w, &codec.MsgpackHandle{}),
		}
		n.connPool[dest] = c
	}

	if err := c.w.WriteByte(rpcType); err != nil {
		return err
	}
	if err := c.encode.Encode(msg); err != nil {
		return err
	}
	if err := c.encode.Encode(sig); err != nil {
		return err
	}
	if err := c.w.Flush(); err != nil {
		return err
	}
	return nil
}

func NewnetworkDealer(filepath string) (*NetworkDealer, error) {
	c := config.Loadconfig(filepath)

	h := makeHost(c.Port, c.Prvkey)
	n := &NetworkDealer{
		connPool:   make(map[string]*conn),
		msgch:      make(chan MsgWithSig, 100),
		H:          h,
		shutdown:   false,
		shutdownCh: make(chan struct{}),

		reflectedTypesMap: make(map[uint8]reflect.Type),
		config:            c,
	}

	return n, nil
}