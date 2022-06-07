package model

import "net"

//服务端返回信息应该选择UserInfo
type User struct {
	UserId       uint64 `json:"userId"` //限制ID长度小于等于16
	UserName     string `json:"userName"`
	UserPwd      string `json:"userPwd"`
	Sex          string `json:"sex"`
	RegisterTime string `json:"registerTime"`
}

//用于与客户端交互,不包含用户密码
type UserInfo struct {
	UserId       uint64 `json:"userId"` //限制ID长度小于等于16
	UserName     string `json:"userName"`
	Sex          string `json:"sex"`
	RegisterTime string `json:"registerTime"`
	Status       string `json:"status"` //用户状态
}

//简单版用户信息
type SimpleUserInfo struct{
	UserId       uint64 `json:"userId"` //限制ID长度小于等于16
	UserName     string `json:"userName"`
}


type UserStatus = string

//用户状态
const (
	Online  UserStatus = "在线"
	Offline UserStatus = "离线"
)

//用于维护在线用户列表(服务端)
type OnlineUserInfo struct {
	Conn net.Conn `json:"conn"`
	Info UserInfo `json:"info"`
}

//向通道CH_OnlineUser发送UpdateOperate指令可更新在线用户列表
//具体操作由函数UpdateOnlineUser执行
type UpdateOperate struct {
	OPcode   UpdateOperateCode
	UserInfo OnlineUserInfo
}

//在线用户列表(服务端)
var OnlineUsers = make(map[uint64]OnlineUserInfo, 100)

//为保证map线程安全，通过对一个全局的channel发送指令来操作map(服务端)
var CH_OnlineUser = make(chan UpdateOperate, 100)

//更新在线用户列表的操作码
type UpdateOperateCode = int

const (
	//移除
	Remove UpdateOperateCode = iota
	//添加
	Add UpdateOperateCode = iota
)
