package core

type Blockchain struct {
	Vertices map[int]*Round
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Vertices: make(map[int]*Round),
	}
}
