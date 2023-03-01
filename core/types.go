package core

import "reflect"

const (
	FMsgTag uint8 = iota
	LMsgTag
	MMsgTag
)

const N = 6

type Froundmsg struct {
	Mroundmsg
}

type Lroundmsg struct {
	Mroundmsg
}

type Mroundmsg struct {
	RN              uint32   `json:"rn"`
	ReferencesHash  [][]byte `json:"referenceshash"`
	ReferencesIndex []uint8  `json:"referencesindex"`
	Source          string   `json:"source"`
	Hash            []byte   `json:hash`
}

type Message interface {
	Encode() ([]byte, error)
	DisplayinJson() error
	GetRefs() ([][]byte, []uint8)
	HavePath(msg Message, msgbyrounds []*MSGByRound) (bool, error)
	IsEqual(msg Message) (bool, error)
	GetRN() uint32
}

var fmsg Froundmsg
var lmsg Lroundmsg
var mmsg Mroundmsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	MMsgTag: reflect.TypeOf(mmsg),
}
