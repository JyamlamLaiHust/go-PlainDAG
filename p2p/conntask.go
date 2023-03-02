package p2p

import (
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func Startpeer(port int, prvkey crypto.PrivKey, reflectedTypesMap map[uint8]reflect.Type) (*NetworkDealer, error) {
	n, err := NewnetworkDealer(port, prvkey, reflectedTypesMap)
	if err != nil {
		return nil, err
	}
	n.Listen()
	return n, nil

}

func (n *NetworkDealer) Connectpeers(peerid int, idaddrmap map[int]string, idportmap map[int]int, pubstringsmap map[int]string) error {

	for id, addr := range idaddrmap {
		if id != peerid {

			writer, err := n.Connect(idportmap[id], addr, pubstringsmap[id])
			if err != nil {
				return err
			}
			log.Println("connect to ", addr, idportmap[id], " success")
			n.connPool[addr+strconv.Itoa(idportmap[id])] = &conn{
				w: writer,

				dest: addr + strconv.Itoa(idportmap[id]),

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
