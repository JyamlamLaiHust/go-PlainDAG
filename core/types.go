package core

import "reflect"

const (
	FMsgTag uint8 = iota
	LMsgTag
	PMsgTag
)

type froundmsg struct {
}

type lroundmsg struct {
}

type plainmsg struct{}

var fmsg froundmsg
var lmsg lroundmsg
var pmsg plainmsg
var ReflectedTypesMap = map[uint8]reflect.Type{
	FMsgTag: reflect.TypeOf(fmsg),
	LMsgTag: reflect.TypeOf(lmsg),
	PMsgTag: reflect.TypeOf(pmsg),
}
