package main

import (
	"bufio"
	"net"

	"github.com/hashicorp/go-msgpack/codec"
)

// NetConn represents a connection established from one node to another.
type NetConn struct {
	target string
	conn   net.Conn
	w      *bufio.Writer
	enc    *codec.Encoder
}

// Release closes the connection in a NetConn variable.
func (n *NetConn) Release() error {
	return n.conn.Close()
}
