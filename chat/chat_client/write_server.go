package main

import (
	"chat/proto"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func WirteServer(conn net.Conn, cmd, data string) {
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
		return
	}

}
