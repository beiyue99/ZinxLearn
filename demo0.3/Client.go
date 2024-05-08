package main

import (
	"fmt"
	"net"
	"time"
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
		//调用write写数据
		_, err := conn.Write([]byte("hello zinxV0.2..."))
		if err != nil {
			fmt.Println("conn Write err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf err")
			return
		}
		fmt.Printf("server callback %s,cnt=%d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}
