package proto 
const (
	UserLogin           = "login"
	UserRegister        = "register"
	UserSendMessageCmd  = "@"
	UserHelp       = "help"
	UserOnline= "getOnline"

)
const Helpstr=`
1.login 命令 登录 例子： login 账号 密码-----login 1 123
2.register 命令 注册 离职 register 账号 密码 名称 性别 ----- register 2 123 yiwei man
3.@ 命令给指定用户发送命令 @账号Id 发送内容----- @1 你好
4.help 命令获取命令说明
5.getOnline 获取在线用户
`