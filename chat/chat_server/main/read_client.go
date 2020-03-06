package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"program/chat/chat_server/common"
	"program/chat/proto"
)

//接收用户发的信息
func ReadClient(conn net.Conn) {
	for {
		buf := make([]byte, 1024*2)
		n, err := conn.Read(buf) //读取内容长度
		if n != 4 {
			fmt.Println("协议错误")
		}
		var packLen uint32
		packLen = binary.BigEndian.Uint32(buf[0:4])
		n, err = conn.Read(buf[0:packLen]) //读取内容
		if n != int(packLen) {
			fmt.Println("协议包长度不对")
			return
		}
		if err != nil {
			fmt.Println("conn.Read failed", err)
			return
		}

		//json 转对象
		var userRequest proto.Message
		json.Unmarshal(buf[0:packLen], &userRequest)
		switch userRequest.Cmd {
		case "register":
			_,Msg, _ := common.Register(userRequest.Data)
			//fmt.Println(Msg)
			WriteClient(Msg,conn)
		case "login":
			UserId,Msg, _ :=common.Login(userRequest.Data)
			//fmt.Println(Msg)
			WriteClient(Msg,conn)
			if UserId != -1{
				  cl  :=&Client{
					conn:conn,
					userId:UserId,
				}
				//添加到client中				
				clientMgr.AddClient(UserId,cl)
				WriteClient("在线人列表:\r\n",conn)
				//显示所用在线的用户
				for _,v :=range clientMgr.GetAllUsers(){
					WriteClient(fmt.Sprintf("用户Id:%d, 用户名称：状态：在线\r\n",v.userId),conn)
				}
			}
		}
	}
}
