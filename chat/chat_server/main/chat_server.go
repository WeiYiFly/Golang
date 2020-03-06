package main

import (
	"net"
	"fmt"
	"encoding/binary"
)

func main () {
	//开启Tcp
	ln,err := net.Listen("tcp","127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err :",err)
		return
	}
	defer ln.Close()
	fmt.Println("==========聊天服务器 已开启 =======")
	for{
		conn ,err :=ln.Accept()
		if err != nil{
			fmt.Println("ln.Accept err : ",err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn){
	defer conn.Close()
	//接收用户发的信息
	for{
		buf := make([]byte,1024*2)	
		n,err:=conn.Read(buf)//读取内容长度
		if n !=4 {
			fmt.Println("read header failed",err)
		}
		var packLen uint32
		packLen = binary.BigEndian.Uint32(buf[0:4])
		n, err = conn.Read(buf[0:packLen]) //读取内容
		if n != int(packLen) {
			fmt.Println("read header failed",err)
			return
		}
		fmt.Println("data is ",string(buf[0:packLen]))

	}
}
