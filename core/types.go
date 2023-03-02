package core

import (
	"reflect"
)

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

// Ref is used to refer a message, and a index field is added to make fast the searching procedure
// Honestly adding this index field is not recommended because not all nodes have the same index-pubkey mapping
type Mroundmsg struct {
	Rn         uint32 `json:"rn"`
	References []Ref  `json:"references"`
	Source     []byte `json:"source"`
	Hash       []byte `json:hash`
}

type Ref struct {
	H     []byte
	Index uint8
}

type Message interface {
	Encode() ([]byte, error)
	DisplayinJson() error
	//MarshalJSON() ([]byte, error)
	GetRefs() []Ref
	HavePath(msg Message, msgbyrounds []*MSGByRound, targetmsground *MSGByRound) (bool, error)

	GetRN() uint32
	GetHash() []byte
	VerifySig(*Node, []byte) (bool, error)
}

var fmsg Froundmsg
var lmsg Lroundmsg
var mmsg Mroundmsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	MMsgTag: reflect.TypeOf(mmsg),
}
