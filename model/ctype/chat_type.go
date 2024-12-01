package ctype

type MsgType int

const (
	InRoomMsg MsgType = 1 + iota
	TextMsg
	ImageMsg
	OutRoomMsg
	SystemMsg
)
