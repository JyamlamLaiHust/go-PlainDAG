package core

type Blockchain struct {
	Vertices map[int]*MSGByRound
}

func NewBlokchain() *Blockchain {
	return &Blockchain{
		Vertices: make(map[int]*MSGByRound),
	}
}
