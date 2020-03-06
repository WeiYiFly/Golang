package main

import (
	"fmt"
	"net"
	"program/chat/redis"
	"time"
)

func main() {
	//开启Tcp
	// ln, err := net.Listen("tcp", "127.0.0.1:8888")
	// if err != nil {
	// 	fmt.Println("net.Listen err :", err)
	// 	return
	// }
	// defer ln.Close()
	err :=runServer("127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}
	fmt.Println("==========聊天服务器 已开启 =======")

	//初始化redis
	redis.InitRedis("localhost:6379", 16, 1024, time.Second*300)

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		fmt.Println("ln.Accept err : ", err)
	// 		continue
	// 	}
	// 	go handleConnection(conn)
	// }

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	//接收用户发的信息
	ReadClient(conn)
}
