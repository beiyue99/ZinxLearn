package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

//封包，拆包具体模块

type DataPack struct {
}

//封包拆包实例的初始化方法

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	return 8 //消息id+消息len
}

// 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将datalen写进dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据写进dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil //最终返回一个二进制文件
}

// 拆包方法  先把Head读出来，然后得到包体长度，再读包体内容
func (dp *DataPack) Unpack(binaryDate []byte) (ziface.IMessage, error) {
	//创建一个读二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryDate)
	//通过Head得到 id和len
	msg := &Message{}
	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断包体是否过长
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
