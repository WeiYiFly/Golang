package main

import (
	"fmt"
	"program/chat/redis"
	"time"
)

func main() {

	fmt.Println("==========聊天服务器 已开启 =======")

	//初始化redis
	redis.InitRedis("localhost:6379", 16, 1024, time.Second*300)

	//tcp 要最后开启
	//err := runServer("127.0.0.1:9000")
	err := runServer("172.16.0.4:8888")
	if err != nil {
		fmt.Println("net.Listen err :", err)
		//return
	}
	for{

	}

}
