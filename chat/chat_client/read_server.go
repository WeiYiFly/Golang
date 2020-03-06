package main

import (
	"fmt"
	"net"
)

func ReadServer(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf) //接收服务器的请求
		if err != nil {
			fmt.Println("conn.Read err: ", err)
			return
		}
		fmt.Println(string(buf[:n])) //打印接收到数据
	}
}
