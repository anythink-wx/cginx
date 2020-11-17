package iface

type Imessage interface {
	GetMsgId() uint16

	GetMsgLen() uint32

	GetMsgData() []byte

	SetMsgId(uint16)
	SetMsgLen(uint32)
	SetMsgData([]byte)
}
