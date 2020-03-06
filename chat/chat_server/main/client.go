package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"program/chat/chat_server/common"
	"program/chat/proto"
	"errors"
)

type Client struct {
	conn   net.Conn
	userId int
	buf    [8192]byte
}


func (p *Client) Process() (err error) {

	for {
		var msg proto.Message
		msg, err = p.readPackage()
		if err != nil {
			//	fmt.Println("")
			clientMgr.DelClient(p.userId)
			//TODO:通知所有在线用户，该用户已经下线
			return err
		}

		//err = p.processMsg(msg)
		if err != nil {
			fmt.Println("process msg failed, err:", err)
			continue
			//return
		}
	}
}
//读包
func (p *Client) readPackage() (msg proto.Message, err error) {
		n, err := p.conn.Read(p.buf[0:4]) //读取内容长度
		if n != 4 {
			fmt.Println("协议错误")
			return
		}
		var packLen uint32
		packLen = binary.BigEndian.Uint32(p.buf[0:4])
		n, err = p.conn.Read(p.buf[0:packLen]) //读取内容
		if n != int(packLen) {
			fmt.Println("协议包长度不对")
			return
		}
		if err != nil {
			fmt.Println("conn.Read failed", err)
			return
		}

		//json 转对象
		//var userRequest proto.Message
		json.Unmarshal(p.buf[0:packLen], &msg)
		switch userRequest.Cmd {
		case "register":
			_,Msg, _ := common.Register(userRequest.Data)
			//fmt.Println(Msg)
			//WriteClient(Msg,conn)
			p.writePackage([]byte(Msg))
		case "login":
			UserId,Msg, _ :=common.Login(userRequest.Data)
			//fmt.Println(Msg)
			//WriteClient(Msg,conn)
			p.writePackage([]byte(Msg))
			/**
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
			}**/
		}
		return
}
//写包
func (p *Client) writePackage(data []byte) (err error) {

	packLen := uint32(len(data))
	binary.BigEndian.PutUint32(p.buf[0:4], packLen)
	n, err := p.conn.Write(p.buf[0:4])
	if err != nil {
		fmt.Println("write data  failed")
		return
	}

	n, err = p.conn.Write(data)
	if err != nil {
		fmt.Println("write data  failed")
		return
	}

	if n != int(packLen) {
		fmt.Println("write data  not finished")
		err = errors.New("write data not fninshed")
		return
	}
	return
}
//消息处理
func (p *Client) processMsg(msg proto.Message) (err error) {

	switch msg.Cmd {
	case proto.UserLogin:
		err = p.login(msg)
	case proto.UserRegister:
		err = p.register(msg)
	case proto.UserSendMessageCmd:
		//err = p.proccessUserSendMessage(msg)
	default:
		err = errors.New("unsupport message")
		return
	}
	return
}