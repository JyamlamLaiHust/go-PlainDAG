package p2p

import (
	"log"

	"github.com/hashicorp/go-msgpack/codec"
)

func Startpeer(filepath string) (*NetworkDealer, error) {
	h, err := NewnetworkDealer(filepath)
	if err != nil {
		return nil, err
	}
	h.Listen()
	return h, nil

}

func (n *NetworkDealer) Connectpeers() {
	c := n.config
	for id, addr := range c.IdaddrMap {
		if id != c.Id {

			writer, err := n.Connect(c.IdportMap[id], addr, c.Pubkeyothersmap[id])
			if err != nil {
				panic(err)
			}
			log.Println("connect to ", addr, c.Port, " success")
			n.connPool[addr] = &conn{
				w:      writer,
				dest:   addr,
				encode: codec.NewEncoder(writer, &codec.MsgpackHandle{}),
			}
		}
	}

}

// func printconfig(c *config.P2pconfig) {
// 	log.Println("id: ", c.Id)
// 	log.Println("nodename: ", c.Nodename)
// 	log.Println("port: ", c.Port)
// 	log.Println("addr: ", c.Ipaddress)
// 	log.Println("idportmap: ", c.IdportMap)
// 	log.Println("idaddrmap: ", c.IdaddrMap)
// 	log.Println("pubkeyothersmap: ", c.Pubkeyothersmap)
// 	log.Println("prvkey: ", c.Prvkey)
// 	log.Println("pubkey: ", c.Pubkey)

// }
