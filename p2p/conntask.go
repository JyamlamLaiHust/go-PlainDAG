package p2p

import (
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/go-msgpack/codec"
)

func Startpeer(filepath string, reflectedTypesMap map[uint8]reflect.Type) (*NetworkDealer, error) {
	n, err := NewnetworkDealer(filepath, reflectedTypesMap)
	if err != nil {
		return nil, err
	}
	n.Listen()
	return n, nil

}

func (n *NetworkDealer) Connectpeers() error {
	c := n.config
	for id, addr := range c.IdaddrMap {
		if id != c.Id {

			writer, err := n.Connect(c.IdportMap[id], addr, c.Pubkeyothersmap[id])
			if err != nil {
				return err
			}
			log.Println("connect to ", addr, c.IdportMap[id], " success")
			n.connPool[addr+strconv.Itoa(c.IdportMap[id])] = &conn{
				w: writer,

				dest: addr + strconv.Itoa(c.IdportMap[id]),

				encode: codec.NewEncoder(writer, &codec.MsgpackHandle{}),
			}
		}
	}
	return nil

}

// func (n *NetworkDealer) PrintConnPool() {
// 	for k, v := range n.connPool {
// 		log.Println(k, v)
// 	}

// }

func (n *NetworkDealer) Broadcast(messagetype uint8, msg interface{}, sig []byte) error {

	for _, conn := range n.connPool {
		//fmt.Println(conn)
		err := n.SendMsg(messagetype, msg, sig, conn.dest)
		if err != nil {
			return err
		}
	}
	return nil

}

// func (n *NetworkDealer) HandleMsgForever() {

// 	for {
// 		select {
// 		case <-n.shutdownCh:
// 			return
// 		case msg := <-n.msgch:
// 			log.Println("receive msg: ", msg.Msg)
// 		}

// 	}
// }

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
