package cnet

type message struct {
	Id      uint16
	DataLen uint32
	Data    []byte
}

func NewMessagePackage(id uint16, data []byte) *message {
	return &message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}

}
func (m *message) GetMsgId() uint16 {
	return m.Id
}

func (m *message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *message) GetMsgData() []byte {
	return m.Data
}

func (m *message) SetMsgId(id uint16) {
	m.Id = id
}
func (m *message) SetMsgLen(ln uint32) {
	m.DataLen = ln
}
func (m *message) SetMsgData(data []byte) {
	m.Data = data
}
