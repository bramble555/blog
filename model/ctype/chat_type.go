package ctype

type MsgType int

const (
	InRoomMsg MsgType = 1 + iota // InRoomMsg 将是 1
	TextMsg                      // TextMsg 将是 2
	ImageMsg                     // ImageMsg 将是 3
	SystemMsg                    // SystemMsg 将是 4
)
