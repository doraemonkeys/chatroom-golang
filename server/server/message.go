package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"

	"github.com/garyburd/redigo/redis"
)

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

func ReceiveMessage(conn net.Conn) (model.Messsage, error) {
	buf := make([]byte, 1024*100)
	var msg model.Messsage
	n, err := conn.Read(buf)
	//如何判断是否与客户端断开连接
	if err != nil && n == 0 {
		fmt.Println(conn.RemoteAddr().String(), "客户端连接断开！")
		fmt.Println("err:", err)
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
	case model.GetUnReadMsgType:
		var NewMsg model.GetUnReadMsg
		err := json.Unmarshal([]byte(msg.Data), &NewMsg)
		if err != nil {
			fmt.Println("in func ParseMessage,json.Unmarshal failed")
			return nil, err
		}
		return NewMsg, nil
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

//封装数据用于发送(空接口版本)
// func MessagePacking(MsgType int, MsgData any) (model.Messsage, error) {
// 	switch MsgType {
// 	case model.GroupChatType:
// 		data, err := json.Marshal(MsgData.(model.GroupChatMsg))
// 		if err != nil {
// 			return model.Messsage{}, err
// 		}
// 		var Msg = model.Messsage{
// 			MsgType: model.GroupChatType,
// 			Data:    string(data),
// 		}
// 		Msg.DataLength = len(Msg.Data)
// 		return Msg, err
// 	case model.PrivateChatType:
// 		data, err := json.Marshal(MsgData.(model.PrivateChatMsg))
// 		if err != nil {
// 			return model.Messsage{}, err
// 		}
// 		var Msg = model.Messsage{
// 			MsgType: model.PrivateChatType,
// 			Data:    string(data),
// 		}
// 		Msg.DataLength = len(Msg.Data)
// 		return Msg, err
// 	default:
// 		return model.Messsage{}, model.Err_UnknowMsg
// 	}
// 	return model.Messsage{}, model.Err_UnknowMsg
// }

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

//检查未读消息，并提醒客户端
func CheckUnResdMsg(conn net.Conn, NewLoginMsg model.LoginMsg) error {
	redisConn := redisPool.Get()
	key := "PTMSG" + utils.Uint64ToStr(NewLoginMsg.UserId)
	idHave, err := redis.Bool(redisConn.Do("exists", key))
	if err != nil {
		fmt.Println("in func CheckUnResdMsg,在redis中查询id失败")
		return err
	}
	if !idHave {
		return nil
	}
	unreadNum, err := redis.Int(redisConn.Do("llen", key))
	if err != nil {
		return err
	}
	msgData := model.UnReadMsgNotification{
		MsgNum: unreadNum,
	}
	msg, err := MessagePacking(model.UnReadMsgNotificationType, msgData)
	if err != nil {
		return err
	}
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}
	return nil
}
