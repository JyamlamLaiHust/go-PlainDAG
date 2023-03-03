package core

import (
	"reflect"
)

const (
	FMsgTag uint8 = iota

	BMsgTag
	LMsgTag
)

const f = 1
const rPerwave = 3

var messageconst []byte

type Froundmsg struct {
	BasicMsg
}

type Lroundmsg struct {
	BasicMsg
}

// Ref is used to refer a message, and a index field is added to make fast the searching procedure
// Honestly adding this index field is not recommended because not all nodes have the same index-pubkey mapping
type BasicMsg struct {
	Rn         int      `json:"rn"`
	References [][]byte `json:"references"`
	Source     []byte   `json:"source"`
	Hash       []byte   `json:hash`
	plaintext  []byte   //`json:plaintext`
}

// type Ref struct {
// 	H     []byte
// 	Index uint8
// }

type Message interface {
	Encode() ([]byte, error)
	DisplayinJson() error
	//MarshalJSON() ([]byte, error)
	GetRefs() [][]byte
	HavePath(msg Message, msgbyrounds []*Round, targetmsground *Round) (bool, error)

	GetRN() int
	GetHash() []byte
	VerifySig(*Node, []byte) (bool, error)
	VerifyFields(*Node) error
}

var fmsg Froundmsg
var lmsg Lroundmsg
var bmsg BasicMsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	BMsgTag: reflect.TypeOf(bmsg),
}
