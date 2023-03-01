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
	Rn         uint32 `json:"rn"`
	References []Ref  `json:"references"`
	Source     string `json:"source"`
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
	SetSource(string)
}

var fmsg Froundmsg
var lmsg Lroundmsg
var mmsg Mroundmsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	MMsgTag: reflect.TypeOf(mmsg),
}
