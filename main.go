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
	time.Sleep(15 * time.Second)
	n.Connectpeers()
	select {}

}
