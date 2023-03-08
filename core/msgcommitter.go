package core

import (
	"fmt"
	"sync"
)

type StaticCommitter struct {
	decidedWaveNumber int

	unCommittedMsgs []*slot
	unCommittedlock sync.RWMutex

	unEmbeddedMsgs []*slot
	unEmbeddedlock sync.RWMutex

	// leaderSlot     map[int]int
	// leaderSlotLock sync.RWMutex

	voteMap  map[int]map[string]int
	voteLock sync.RWMutex
	n        *Node
}

type slot struct {
	slotindex int
	msgs      []Message
}

func NewStaticCommitter(n *Node) *StaticCommitter {
	return &StaticCommitter{
		decidedWaveNumber: -1,
		unCommittedMsgs:   make([]*slot, 0),
		unEmbeddedMsgs:    make([]*slot, 0),
		voteMap:           make(map[int]map[string]int),
		//leaderSlot:        make(map[int]int),
		n: n,
	}
}

func (c *StaticCommitter) addVote(waveNumber int, target string) {
	c.voteLock.Lock()
	defer c.voteLock.Unlock()
	if _, ok := c.voteMap[waveNumber]; !ok {
		c.voteMap[waveNumber] = make(map[string]int)
	}
	c.voteMap[waveNumber][target]++
	fmt.Println("The leader in wave     " + fmt.Sprint(waveNumber) + "   has got " + fmt.Sprint(c.voteMap[waveNumber][target]) + " votes")
}
