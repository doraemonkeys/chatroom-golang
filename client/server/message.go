package server

import (
	"client/model"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

//向服务器发送消息
func SendMessage(conn net.Conn, msg model.Messsage) error {
	if msg.DataLength != len(msg.Data) {
		return model.Err_ErrorMessage
	}
	msgData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("in func SendMessage=,json.Marshal failed")
		return err
	}
	_, err = conn.Write(msgData)
	if err != nil {
		fmt.Println("in func SendMessage,conn.Write failed")
		return err
	}
	return nil
}

//为conn接收一个服务端发来的消息
func ReceiveMessage(conn net.Conn) (model.Messsage, error) {
	buf := make([]byte, 1024*100)
	var msg model.Messsage
	n, err := conn.Read(buf)
	//如何判断是否与服务器断开连接
	if err != nil && n == 0 {
		//fmt.Println(conn.RemoteAddr().String(), "与服务端连接断开！",err)
		return msg, model.Err_ConnectionDown
	}
	if err != nil {
		fmt.Println("in func ReceiveMessage,conn.Read failed,err:", err)
		return msg, err
	}
	err = json.Unmarshal(buf[:n], &msg)
	if err != nil {
		fmt.Println("in func ReceiveMessage,json.Unmarshal failed")
		return msg, err
	}
	if len(msg.Data) != msg.DataLength {
		return msg, model.Err_DataCorruption
	}
	return msg, nil
}

//用于客户端发送请求后接收服务器是否返回消息，
//客户端接收到响应会将值更改为对应的响应类型，
//可以优化成一个全局的队列(用来缓存服务器发送的消息信息)
var MsgResType int = -1

//查看服务端是否响应请求(逻辑不严谨)
func CheckResponse(expectType int) bool {
	start := time.Now()
	for {
		if MsgResType == expectType {
			MsgResType = -1
			return true
		}
		time.Sleep(time.Millisecond * 800)
		now := time.Now()
		//1s内未收到响应则判断为超时
		if now.Sub(start) > time.Second {
			MsgResType = -1
			return false
		}
	}
}

//解析消息，返回对应的结构体数据(可优化成泛型函数)
//未完成
func ParseMessage(msg model.Messsage) (any, error) {
	switch msg.MsgType {
	case model.LoginMsgType:
		var NewLoginMsg model.LoginMsg
		err := json.Unmarshal([]byte(msg.Data), &NewLoginMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewLoginMsg, nil
	case model.RegisterMsgResType:
		var NewRegisterMsgRes model.RegisterMsgRes
		err := json.Unmarshal([]byte(msg.Data), &NewRegisterMsgRes)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewRegisterMsgRes, nil
	case model.RegisterMsgType:
		var NewRegisterMsg model.RegisterMsg
		err := json.Unmarshal([]byte(msg.Data), &NewRegisterMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewRegisterMsg, nil
	case model.LoginMsgResType:
		var NewLoginMsgRes model.LoginMsgRes
		err := json.Unmarshal([]byte(msg.Data), &NewLoginMsgRes)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewLoginMsgRes, nil
	case model.GetUsreInfoType:
		var NewGetUsreInfoMsg model.GetUsreInfoMsg
		err := json.Unmarshal([]byte(msg.Data), &NewGetUsreInfoMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewGetUsreInfoMsg, nil
	case model.GetUsreInfoResType:
		var GetUsreInfoResMsg model.GetUsreInfoResMsg
		err := json.Unmarshal([]byte(msg.Data), &GetUsreInfoResMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return GetUsreInfoResMsg, nil
	case model.GetOnlineUsersType:
		return nil, nil
	case model.OnlineUsersResType:
		var OnlineUsersResMsg model.OnlineUsersResMsg
		err := json.Unmarshal([]byte(msg.Data), &OnlineUsersResMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return OnlineUsersResMsg, nil
	case model.GroupChatType:
		var GroupChatMsg model.GroupChatMsg
		err := json.Unmarshal([]byte(msg.Data), &GroupChatMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return GroupChatMsg, nil
	case model.PrivateChatType:
		var PrivateChatMsg model.PrivateChatMsg
		err := json.Unmarshal([]byte(msg.Data), &PrivateChatMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return PrivateChatMsg, nil
	case model.UnReadMsgNotificationType:
		var UnReadMsg model.UnReadMsgNotification
		err := json.Unmarshal([]byte(msg.Data), &UnReadMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return UnReadMsg, nil
	case model.UnReadMsgResType:
		var NewMsg model.UnReadMsgRes
		err := json.Unmarshal([]byte(msg.Data), &NewMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewMsg, nil
	case model.GetGroupChatType:
		var NewMsg model.GetGroupChatMsg
		err := json.Unmarshal([]byte(msg.Data), &NewMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewMsg, nil
	case model.GetGroupChatResType:
		var NewMsg model.GetGroupChatResMsg
		err := json.Unmarshal([]byte(msg.Data), &NewMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewMsg, nil
	default:
		return nil, model.Err_UnknowMsg
	}
}

//封装数据用于发送(泛型版本)
func MessagePacking[T model.MsgData](MsgType int, MsgData T) (model.Messsage, error) {
	data, err := json.Marshal(MsgData)
	if err != nil {
		return model.Messsage{}, err
	}
	var Msg = model.Messsage{
		MsgType: MsgType,
		Data:    string(data),
	}
	Msg.DataLength = len(Msg.Data)
	return Msg, err
}
