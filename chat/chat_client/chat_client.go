package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()

	//向服务器发送 消息
	go func() {
		// 向服务器 发送键盘输入的内容
		var cmd, data string
		for {
			fmt.Scanf("%d %s\n", &cmd, &data)
			if err != nil {
				fmt.Println("conn.Read err: ", err)
				continue
			}
			//conn.Write(str[:n])
			WirteServer(conn, cmd, data)
		}
	}()
	//接收服务器的消息
	ReadServer(conn)

}
