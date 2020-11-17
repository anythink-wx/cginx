package iface

type Idatapack interface {
	GetHeadLen() uint32

	Pack(msg Imessage) ([]byte, error)
	Unpack([]byte) (Imessage, error)
}
