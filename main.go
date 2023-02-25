package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/PlainDAG/go-PlainDAG/p2p"
)

func main() {
	index := flag.Int("f", 0, "config file path")
	flag.Parse()
	//convert int to string
	filepath := "node" + strconv.Itoa(*index)

	n, err := p2p.Startpeer(filepath)
	if err != nil {
		panic(err)
	}
	//fmt.Println(n.H.ID().Pretty())
	time.Sleep(10 * time.Second)
	n.Connectpeers()
	go n.Broadcast()
	go n.HandleMsgForever()
	//n.PrintConnPool()
	select {}

	// var x interface{}
	// x = 0
	// // var a interface{}
	// // a = &x
	// switch i := x.(type) {
	// case int:
	// 	fmt.Println("int", i)
	// default:
	// 	fmt.Println("none")
	// }
}
