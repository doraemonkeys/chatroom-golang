package model

import (
	"errors"
)

//定义消息类型
const (
	//登录消息
	LoginMsgType = iota
	//登录返回消息
	LoginMsgResType = iota
	//注册消息
	RegisterMsgType = iota
	//注册返回消息
	RegisterMsgResType = iota
	//获取某个用户信息
	GetUsreInfoType = iota
	//返回某个用户信息
	GetUsreInfoResType = iota
	//私聊消息
	PrivateChatType = iota
	//群聊消息
	GroupChatType = iota
	//查询群聊消息记录
	GetGroupChatType = iota
	//响应查询群聊消息记录
	GetGroupChatResType = iota
	//获取在线用户列表
	GetOnlineUsersType = iota
	//返回在线用户列表
	OnlineUsersResType = iota
	//获取未读信息
	GetUnReadMsgType = iota
	//返回未读信息
	UnReadMsgResType = iota
	//未读消息通知
	UnReadMsgNotificationType = iota
)

//定义状态码
const (
	//登录成功
	LoginSucess = 200
	//登录失败
	LoginFailed = 200 + iota
	//注册成功
	RegisterSucess = 200 + iota
	//注册失败
	RegisterFailed = 200 + iota
	//接收到的信息损坏
	TransferError = 200 + iota
	//未知错误
	UnknownError = 200 + iota
)

var Err_ConnectionDown error = errors.New("tcp连接断开")                 //tcp连接断开
var Err_FrequentlyLogin error = errors.New("尝试登录次数过多")               //尝试登录次数过多
var Err_PasswordWrong error = errors.New("密码错误")                     //密码错误
var Err_UserNotExist error = errors.New("找不到用户")                     //找不到用户
var Err_UserAlreadyExist error = errors.New("用户已存在")                 //用户已存在
var Err_DataCorruption error = errors.New("接收到损坏的数据")                //接收到损坏的数据
var Err_ErrorMessage error = errors.New("Messsage长度验证失败")            //Messsage长度验证失败
var Err_InternalServerError error = errors.New("服务器出现了一点点小问题，请稍候再试") //服务器内部错误
var Err_UnknowMsg error = errors.New("不支持的消息类型")                     //不支持的消息类型
var Err_RepeatLogin error = errors.New("用户已经登录")                     //用户已经登录
var Err_NoOnlineUsers error = errors.New("无用户在线")                    //无用户在线
var Err_NoUnReadMsg error = errors.New("没有未读消息")                     //没有未读消息
var Err_UserNotOnline error = errors.New("用户不在线")                    //用户不在线
var Err_GetUnknowMsg error = errors.New("可能接收到其他goroutine所需的信息,请检查") //可能接收到其他goroutine所需的信息,请检查

type Messsage struct {
	MsgType    int    `json:"msgType"`    //消息类型
	Data       string `json:"data"`       //数据
	DataLength int    `json:"dataLength"` //消息长度，用于校验接收到的数据是否有误
}

type RegisterMsg struct {
	User User
}
type RegisterMsgRes struct {
	//状态码
	StatusCode int    `json:"StatusCode"`
	UserId     uint64 `json:"userId"`
	Error      string `json:"error"`
}

type LoginMsg struct {
	UserId  uint64 `json:"userId"`
	UserPwd string `json:"userPwd"`
}

type LoginMsgRes struct {
	//状态码
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type GetUsreInfoMsg struct {
	UserId uint64 `json:"userId"`
}

type GetUsreInfoResMsg struct {
	UserInfos UserInfo `json:"userInfos"`
	Error     string   `json:"error"`
}

type OnlineUsersResMsg struct {
	Users []string `json:"users"` //目前限制返回最多10个用户
	Error string   `json:"error"`
}

type PrivateChatMsg struct {
	SrcUser  SimpleUserInfo `json:"srcUser"`
	DstUser  uint64         `json:"dstUser"`
	ChatTime string         `json:"chatTime"` //"2006.01.02 15:04"
	ChatData string         `json:"chatData"`
	//ChatType int//便于扩展功能
}

type GroupChatMsg struct {
	SrcUser  SimpleUserInfo `json:"srcUser"`
	ChatTime string         `json:"chatTime"` //"2006.01.02 15:04"
	ChatData string         `json:"chatData"`
	//ChatType int//便于扩展功能
}

type GetGroupChatMsg struct {
	SrcUser SimpleUserInfo `json:"srcUser"`
	//GroupNum uint64
}

type GetGroupChatResMsg struct {
	Data  []GroupChatMsg `json:"data"` //目前限制返回最多10条记录
	Error string         `json:"error"`
}

type GetUnReadMsg struct {
	ID uint64 `json:"iD"`
}

type UnReadMsgRes struct {
	SrcUsers []SimpleUserInfo `json:"srcUser"`
	Data     []PrivateChatMsg `json:"data"`
	Error    string           `json:"error"`
}

type UnReadMsgNotification struct {
	SrcUsers []SimpleUserInfo `json:"srcUser"`
	MsgNum   int              `json:"msgNum"`
	Error    string           `json:"error"`
}

type MsgData interface {
	GroupChatMsg | PrivateChatMsg | UnReadMsgNotification | UnReadMsgRes | GetUnReadMsg |
	GetGroupChatMsg | GetGroupChatResMsg
}

