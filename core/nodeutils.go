package core

import (
	"github.com/PlainDAG/go-PlainDAG/config"
	"github.com/PlainDAG/go-PlainDAG/p2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func (n *Node) constructpubkeyMap() error {
	peeridSlice := n.network.H.Peerstore().Peers()

	stringpubkeymap := make(map[string]crypto.PubKey, len(peeridSlice))
	stringidmap := make(map[string]int, len(peeridSlice))
	for _, peerid := range peeridSlice {
		pubkey := n.network.H.Peerstore().PubKey(peerid)
		str := peerid.Pretty()
		pubkeybytes, err := crypto.MarshalPublicKey(pubkey)
		if err != nil {
			return err
		}
		stringpubkeymap[string(pubkeybytes)] = pubkey
		stringidmap[string(pubkeybytes)] = n.cfg.PubkeyIdMap[str]
	}
	n.cfg.StringpubkeyMap = stringpubkeymap
	n.cfg.StringIdMap = stringidmap
	return nil
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
	c.Pubkey = n.H.Peerstore().PubKey(n.H.ID())
	Pubkeyraw, err := crypto.MarshalPublicKey(c.Pubkey)
	if err != nil {
		return nil, err
	}
	c.Pubkeyraw = Pubkeyraw

	node.cfg = c
	node.handler = NewStatichandler(&node)
	node.currentround.Store(0)

	return &node, err
}
