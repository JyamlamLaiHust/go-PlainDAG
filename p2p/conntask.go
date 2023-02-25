package p2p

import (
	"encoding/binary"
	"log"
	"strconv"
	"time"

	"github.com/PlainDAG/go-PlainDAG/core"
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
			log.Println("connect to ", addr, c.IdportMap[id], " success")
			n.connPool[addr+strconv.Itoa(c.IdportMap[id])] = &conn{
				w: writer,

				dest: addr + strconv.Itoa(c.IdportMap[id]),
				//convert int to strin
				encode: codec.NewEncoder(writer, &codec.MsgpackHandle{}),
			}
		}
	}

}

// func (n *NetworkDealer) PrintConnPool() {
// 	for k, v := range n.connPool {
// 		log.Println(k, v)
// 	}

// }

func (n *NetworkDealer) Broadcast() {
	var i uint32
	i = 0
	for {
		time.Sleep(5 * time.Second)
		for _, conn := range n.connPool {

			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, uint32(i))
			i++
			// concatenate bytes

			b := append([]byte("hello 12321312"), bytes...)
			a := core.TestMsg{A: b}
			// serialize a with marshall
			// c, err := json.Marshal(a)
			// if err != nil {
			// 	panic(err)
			// }
			n.SendMsg(0, a, []byte("hello"), conn.dest)
		}
	}

}

func (n *NetworkDealer) HandleMsgForever() {

	for {
		select {
		case <-n.shutdownCh:
			return
		case msg := <-n.msgch:

			switch msgasserted := msg.Msg.(type) {
			case *core.TestMsg:
				log.Println("receive msg: ", msgasserted)
			default:
				log.Println("unknown type of msg")
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
