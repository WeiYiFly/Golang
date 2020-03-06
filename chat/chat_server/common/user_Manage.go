package common

import (
	"strings"
	"strconv"
	"program/chat/chat_server/model"
	"program/chat/redis"
	"encoding/json"
	"fmt"
	git_redis "github.com/garyburd/redigo/redis"
)

func Register(strdata string)(IsSuccess bool, Msg string,errR error){
	Msg=""
	IsSuccess=false
	strsplie:= strings.Split(strdata," ")
	var user model.User
	for i,v :=range strsplie{
		switch i {
		case 0:
			tmp,err:=strconv.Atoi(v)
			if err !=nil {
				Msg="账号 必须是数字"
				errR=err
			}
			user.UserId=tmp
		case 1:
			user.PassWd=v
		case 2:
			user.Nick=v
		case 3:
			user.Sex=v
		case 4:
			user.Header=v
		case 5:
			user.Sex=v
		case 6:
			user.Header=v	
		}
	}
	//转json
	datajson, err := json.Marshal(user) //数据对象转json
	if err!=nil{
		Msg +=fmt.Sprintf(" json.Marshal err:",err)
		errR=err
	}
	//添加到 redis 中
	redisConn := redis.GetConn()
	_, err = redisConn.Do("HSet", model.UserTable, fmt.Sprintf("%d", user.UserId), datajson)
	if Msg== ""{
		Msg ="注册成功"
		IsSuccess=true
	}
	return
}
func Login(strdata string)(UserIdR int,Msg string,errR error){
	Msg=""
	UserIdR=-1
	strsplie:= strings.Split(strdata," ")
	if len(strsplie)!=2{
		Msg="输入的格式不对 密码不能有空格"
		return
	}
	//获取redis连接
	redisConn := redis.GetConn()
	result, err := git_redis.String(redisConn.Do("HGet", model.UserTable, strsplie[0]))
	if err !=nil {
		//Msg =fmt.Sprint("redisConn.Do HGet err:" ,err)
		Msg ="无此用户"
		return
	}
	//转对象
	var user model.User

	json.Unmarshal([]byte(result), &user)
	tmp :=strings.Replace(strsplie[1],"\r","",-1)
	tmp =strings.Replace(tmp,"\n","",-1)
	if user.PassWd != tmp {
		Msg ="密码错误"
		return
	}
	fmt.Printf(result)
	if Msg== ""{
		Msg ="登录成功"
		UserIdR=user.UserId
	}
	return 
}
