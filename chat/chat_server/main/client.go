package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	git_redis "github.com/garyburd/redigo/redis"
	"net"
	"program/chat/chat_server/model"
	"program/chat/proto"
	"program/chat/redis"
	"strconv"
	"strings"
)

type Client struct {
	conn net.Conn
	model.User
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
			p.BroadcastSend(fmt.Sprintf("\r\n用户Id:%d  用户：%s 离线--\r\n", p.userId, p.User.Nick))
			return err
		}
		fmt.Println(msg)
		err = p.processMsg(msg)
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
	json.Unmarshal(p.buf[0:packLen], &msg)
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
		err := p.Login(msg)
		if err != nil {
			fmt.Println("Login err ", err)
		}
	case proto.UserRegister:
		err := p.Register(msg)
		if err != nil {
			fmt.Println("Register err ", err)
		}
	case proto.UserSendMessageCmd:
		err := p.SendMsg(msg)
		if err != nil {
			fmt.Println("SendMsg err ", err)
		}
	case proto.UserHelp:
		p.Help()
	case proto.UserOnline:
		p.GenOnline()
	default:
		err = errors.New("unsupport message")
		p.writePackage([]byte("命令不对"))
		return
	}
	return
}

//登录
func (p *Client) Login(msg proto.Message) (errR error) {
	Msg := ""
	fmt.Println(" userId ", p.userId)
	if p.userId != 0 {
		Msg = "只能登录一个账号"
		p.writePackage([]byte(Msg))
		return
	}
	strdata := msg.Data
	strsplie := strings.Split(strdata, " ")
	if len(strsplie) != 2 {
		Msg = "输入的格式不对 密码不能有空格"
	}
	//获取redis连接
	redisConn := redis.GetConn()
	result, err := git_redis.String(redisConn.Do("HGet", model.UserTable, strsplie[0]))
	if err != nil {
		//Msg =fmt.Sprint("redisConn.Do HGet err:" ,err)
		Msg = "无此用户"
		p.writePackage([]byte(Msg))
		return

	}
	//转对象
	var user model.User

	json.Unmarshal([]byte(result), &user)
	tmp := strings.Replace(strsplie[1], "\r", "", -1)
	tmp = strings.Replace(tmp, "\n", "", -1)
	if user.PassWd != tmp {
		Msg = "密码错误"
		p.writePackage([]byte(Msg))
		return
	}

	if Msg == "" {
		Msg = "登录成功\n"
		_, err1 := clientMgr.GetClient(user.UserId)
		if err1 == nil {
			Msg = "此用户已经登录过了"
			p.writePackage([]byte(Msg))
			return
		}
		cl := &Client{
			conn:   p.conn,
			User:   user,
			userId: user.UserId,
		}
		p.userId = user.UserId
		p.User = user
		//广播所有人 我上线了
		p.BroadcastSend(fmt.Sprintf("\r\n用户Id:%d  用户：%s 上线了--\r\n", p.userId, p.User.Nick))
		//添加到client中
		clientMgr.AddClient(user.UserId, cl)
		//显示所用在线的用户
		tmponlineLise := "在线人列表:\r\n"
		for _, v := range clientMgr.GetAllUsers() {
			tmponlineLise += fmt.Sprintf("用户Id:%d, 用户名称：%s 状态：在线\r\n", v.userId, v.Nick)
		}
		Msg += tmponlineLise

	}
	p.writePackage([]byte(Msg))
	return
}

//注册
func (p *Client) Register(msg proto.Message) (errR error) {
	Msg := ""
	strdata := msg.Data
	strsplie := strings.Split(strdata, " ")
	var user model.User
	for i, v := range strsplie {
		switch i {
		case 0:
			tmp, err := strconv.Atoi(v)
			if err != nil {
				Msg = "账号 必须是数字"
				errR = err
			}
			user.UserId = tmp
		case 1:
			user.PassWd = v
		case 2:
			user.Nick = v
		case 3:
			user.Sex = v
		case 4:
			user.Header = v
		case 5:
			user.Sex = v
		case 6:
			user.Header = v
		}
	}
	//转json
	datajson, err := json.Marshal(user) //数据对象转json
	if err != nil {
		Msg += fmt.Sprint("json.Marshal err:", err)
		errR = err
	}
	//添加到 redis 中
	redisConn := redis.GetConn()
	_, err = redisConn.Do("HSet", model.UserTable, fmt.Sprintf("%d", user.UserId), datajson)
	if Msg == "" {
		Msg = "注册成功"
	}
	p.writePackage([]byte(Msg))
	return
}

//发送信息
func (p *Client) SendMsg(msg proto.Message) (errR error) {
	strdata := msg.Data
	dataSplit := strings.Split(strdata, " ")
	userId := dataSplit[0]
	userSendMsg := strings.Replace(strdata, userId, "", 1)
	userSendMsg = strings.Trim(userSendMsg, " ")
	//寻找需要发送用户的客户连接
	id, _ := strconv.Atoi(userId)
	clientv, err := clientMgr.GetClient(id)
	if err != nil {
		errR = err
		p.writePackage([]byte("没有此用户"))
		return
	}
	tmp := fmt.Sprintf("\r\n用户Id: %d 用户名称：%s 发消息给你内容如下：\r\n", p.userId, p.User.Nick)
	userSendMsg = tmp + "=======:" + userSendMsg
	clientv.writePackage([]byte(userSendMsg))
	p.writePackage([]byte("发送成功"))
	return
}

//获取在线用户列表
func (p *Client) GenOnline() {
	//显示所用在线的用户
	tmponlineLise := "在线人列表:\r\n"
	for _, v := range clientMgr.GetAllUsers() {
		tmponlineLise += fmt.Sprintf("用户Id:%d, 用户名称：%s 状态：在线\r\n", v.userId, v.Nick)
	}
	p.writePackage([]byte(tmponlineLise))
}

//获取帮助
func (p *Client) Help() {
	p.writePackage([]byte(proto.Helpstr))
}

//广播发送
func (p *Client) BroadcastSend(strmsg string) {
	for _, v := range clientMgr.GetAllUsers() {
		v.writePackage([]byte(strmsg))
	}

}
