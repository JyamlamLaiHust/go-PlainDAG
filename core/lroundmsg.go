package core

func NewLroundMsg(arefs [][]byte, b *BasicMsg) (*LRoundMsg, error) {

	return &LRoundMsg{
		BasicMsg: b,
		ARefs:    arefs,
	}, nil
}
