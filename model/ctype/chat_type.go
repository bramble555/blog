package ctype

type MsgType int

const (
	InRoomMsg MsgType = 1 + iota
	OutRoomMsg
	TextMsg
	ImageMsg
	SystemMsg
)
