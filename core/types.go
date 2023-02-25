package core

import "reflect"

type TestMsg struct {
	A []byte
}

type Wrongmsg struct {
	a string
}

const (
	TestMsgTag uint8 = iota
)

var testmsg TestMsg
var wrongmsg Wrongmsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	TestMsgTag: reflect.TypeOf(testmsg),
}
