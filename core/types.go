package core

import (
	"reflect"
)

const (
	FMsgTag uint8 = iota

	BMsgTag
	LMsgTag
	TMsgTag
)

const f = 1
const rPerwave = 3

const Batchsize = 1000

var messageconst []byte

const plainMsgSize = 20

type FRoundMsg struct {
	*BasicMsg
}

type LRoundMsg struct {
	*BasicMsg
	ARefs [][]byte `json:arefs`
}

type BasicMsg struct {
	Rn         int        `json:"rn"`
	References [][]byte   `json:"references"`
	Source     []byte     `json:"source"`
	Hash       []byte     `json:hash`
	plainmsg   []PlainMsg //`json:plaintext`
}

type PlainMsg struct {
	Msg []byte
}

type ThresSigMsg struct {
	Sig []byte `json:sig`
	//wave number
	Wn     int    `json:wn`
	Source []byte `json:source`
}

type Message interface {
	Encode() ([]byte, error)
	DisplayinJson() error
	//MarshalJSON() ([]byte, error)
	GetRefs() [][]byte
	HavePath(msg Message, msgbyrounds []*Round, targetmsground *Round) (bool, error)

	GetRN() int
	GetHash() []byte
	GetSource() []byte

	VerifyFields(*Node) error
}

var fmsg FRoundMsg
var lmsg LRoundMsg
var bmsg BasicMsg
var tmsg ThresSigMsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	BMsgTag: reflect.TypeOf(bmsg),
	TMsgTag: reflect.TypeOf(tmsg),
}
