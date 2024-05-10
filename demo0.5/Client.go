package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

// 模拟客户端

func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinxV0.5 client Test Message")))
		if err != nil {
			fmt.Println("pack err ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write err", err)
			return
		}
		//读取服务器回过来的消息,先读取包头
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err", err)
			break
		}
		//把读到的head二进制数据组装成message类型
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client Unpack err", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//第二次读取，把data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data err", err)
				return
			}
			fmt.Println("recv Server msg,Id=", msg.Id, "len=", msg.DataLen, "data=", string(msg.Data))

		}

		time.Sleep(time.Second)
	}
}
