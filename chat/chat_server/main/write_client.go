package main

import (
	"fmt"
	"net"
)

//接收用户发的信息
func WriteClient(Msg string,conn net.Conn) {
	//发送内容
	_, err := conn.Write([]byte(Msg))
	if err != nil {
		fmt.Println("conn.Write err ",err)
	}
}
