package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {

	// 模拟服务器
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listenner err!", err)
		return
	}

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("Server Accept err", err)
			}

			go func(conn net.Conn) {
				//拆包的过程
				dp := NewDataPack()
				for {
					//先读head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server Unpack err", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//有数据,开始读取
						//因为msgHead是抽象指针，无法.出属性，需要转化为实例
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据datalen的长度从io流读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server Unpack data err", err)
							return
						}
						//完整的一个数据包读取完毕
						t.Log("Recv MsgID:", msg.Id, "datalen=", msg.DataLen, "data=", string(msg.Data))
					}
				}
			}(conn)

		}
	}()

	//模拟客户端

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}

	dp := NewDataPack()

	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	//封装数据包
	snedData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err", err)
		return
	}
	snedData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err", err)
		return
	}
	//将两个包粘在一起
	snedData1 = append(snedData1, snedData2...)

	if _, err := conn.Write(snedData1); err != nil {
		fmt.Println("send err", err)
		return
	}

	time.Sleep(time.Second)
}
