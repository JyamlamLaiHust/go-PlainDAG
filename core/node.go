package core

import "sync/atomic"

type Node struct {
	//DAG ledger structure
	bc *Blockchain `json:"bc"`

	//thread-safe integer
	currentround atomic.Uint32 `json:"currentround"`
}

func (n *Node) GenTrans() {

}
