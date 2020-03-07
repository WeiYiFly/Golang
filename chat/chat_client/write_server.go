package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"program/chat/proto"
	"strings"
)

//向服务器发送命令
func WirteServer(conn net.Conn, strdata string) {
	//fmt.Println("开始===发送")
	//fmt.Println("strdata = ",strdata)
	var cmd, data string
	strsplie := strings.Split(strdata, " ")
	if len(strsplie) < 1 {
		fmt.Println("请输入正确的命令")
		return
	}
	cmd = strsplie[0]
	data = strings.Replace(strdata, strsplie[0], "", 1)
	data = strings.Trim(data, " ")
	if data == "" {
		cmd = strings.Replace(cmd, "\r", "", 1)
		cmd = strings.Replace(cmd, "\n", "", 1)
	}
	if strings.HasPrefix(strdata, "@") {
		cmd = "@"
		data = strings.Replace(strdata, "@", "", 1)
	}

	//组包
	var CmdMessage proto.Message
	CmdMessage.Cmd = cmd
	CmdMessage.Data = data
	datajson, err := json.Marshal(CmdMessage) //数据对象转json

	//前面四位是表示内容长度
	var buf [4]byte
	packLen := uint32(len(datajson))
	binary.BigEndian.PutUint32(buf[0:4], packLen)
	n, err := conn.Write(buf[:4])
	if err != nil || n != 4 {
		fmt.Println("write data failed")
	}
	//发送内容
	_, err = conn.Write(datajson)
	if err != nil {
		fmt.Println("write data failed")
		return
	}

}
