package main

import (
	"fmt"
	"net"
	"os"
	"program/chat/proto"
)

func main() {
	conn, err := net.Dial("tcp", "182.61.14.181:8888")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("===========================简单版聊天客户端======================")
	fmt.Println(proto.Helpstr)
	fmt.Println("前先登录")
	fmt.Println("===============================================================")
	fmt.Printf("########:")
	//向服务器发送 消息
	go func() {
		//从键盘输入内容，给服务器发送内容
		str := make([]byte, 1024)
		for {

			n, err := os.Stdin.Read(str) //从键盘读取内容，放在str
			if err != nil {
				fmt.Println("OS.Stdin err；", err)
				return
			}
			//conn.Write(str[:n])
			WirteServer(conn, string(str[:n]))
		}
	}()

	//接收服务器的消息
	ReadServer(conn)

}
