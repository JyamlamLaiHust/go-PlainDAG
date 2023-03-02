package core

type Blockchain struct {
	Vertices map[int]*MSGByRound
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Vertices: make(map[int]*MSGByRound),
	}
}
