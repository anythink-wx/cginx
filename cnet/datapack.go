package cnet

import (
	"bytes"
	"cginx/iface"
	"cginx/utils"
	"encoding/binary"
	"errors"
	"fmt"
)

type DataPack struct {
}

func NewDataPack() *DataPack {

	return &DataPack{}

}

//封包方法 len|msgId|data
func (d *DataPack) Pack(msg iface.Imessage) (buffer []byte, err error) {
	buf := bytes.NewBuffer([]byte{})

	err = binary.Write(buf, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		return nil, err
	}

	buffer = buf.Bytes()
	return
}

//拆包方法,读 包的head 信息， 根据head信息，读取后面长度
func (d *DataPack) Unpack(b []byte) (imessage iface.Imessage, err error) {

	msg := &message{}

	//创建 字节读取器
	reader := bytes.NewReader(b)

	//读取包头
	err = binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	//读取msgid
	err = binary.Read(reader, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}



	//判断datalen 是否超过配置大小
	if utils.ServerOpt.MaxPackageSize > 0 && msg.DataLen > utils.ServerOpt.MaxPackageSize {
		mLen := fmt.Sprintf("%d",msg.DataLen)
		return nil, errors.New("package len > cginx.json MaxPackageSize ,current package is " + mLen)
	}

	return msg, err

}

func (d *DataPack) GetHeadLen() uint32 {
	//4 字节 包体大小， 2字节 Msgid
	return 6
}
