package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"

	"github.com/garyburd/redigo/redis"
)

func ProcessPrivateChat(conn net.Conn, msgData model.PrivateChatMsg) {
	redisConn := redisPool.Get()
	_, err := redis.String(redisConn.Do("hget", "users", msgData.DstUser))
	if err != nil {
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送私聊聊消息失败!","找不到ID:",msgData.DstUser)
		ResponseGetUsreInfo(conn,model.UserInfo{},model.Err_UserNotExist)
		return 
	}
	//如果私聊对象不在线
	v, ok := model.OnlineUsers[msgData.DstUser]
	if !ok {
		data, err := json.Marshal(msgData)
		if err != nil {
			fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送私聊聊消息失败")
			return
		}
		redisConn := redisPool.Get()
		defer redisConn.Close()
		key := "PTMSG" + utils.Uint64ToStr(msgData.DstUser)
		redisConn.Do("rpush", key, string(data))
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送私聊聊消息", "但对方ID:", msgData.DstUser, "不在线,无法转发私聊消息，已存入数据库！")
		return
	}
	msg, err := MessagePacking(model.PrivateChatType, msgData)
	if err != nil {
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送私聊消息失败", err)
		return
	}
	err = SendMessage(v.Conn, msg)
	if err != nil {
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送私聊消息失败", err)
		return
	}
	fmt.Println("用户ID:", msgData.SrcUser.UserId, "发送了一条私聊消息给用户ID:", msgData.DstUser)
	fmt.Printf("内容为:[%s]\n", msgData.ChatData)
}

func ProcessGetUnReadPvMsg(conn net.Conn, msg model.GetUnReadMsg) {
	fmt.Println("用户ID:",msg.ID,"请求查询未读私聊消息")
	redisConn := redisPool.Get()
	key := "PTMSG" + utils.Uint64ToStr(msg.ID)
	idHave, err := redis.Bool(redisConn.Do("exists", key))
	var resMsg = model.UnReadMsgRes{
		Error: "",
	}
	if err != nil {
		fmt.Println("in func ProcessGetUnReadPvMsg,在redis中查询id失败!", err)
		ResponseGetUnReadPvMsg(conn, msg.ID, model.UnReadMsgRes{Error: model.Err_InternalServerError.Error()})
		return
	}
	if !idHave {
		ResponseGetUnReadPvMsg(conn, msg.ID, model.UnReadMsgRes{Error: model.Err_NoUnReadMsg.Error()})
		return
	}
	data, err := redis.Strings(redisConn.Do("lrange", key, 0, -1))
	if err != nil {
		fmt.Println("in func ProcessGetUnReadPvMsg,在redis中查询id失败!", err)
		ResponseGetUnReadPvMsg(conn, msg.ID, model.UnReadMsgRes{Error: model.Err_InternalServerError.Error()})
		return
	}
	PrivateChat := make([]model.PrivateChatMsg,len(data))
	for k, v := range data {
		err := json.Unmarshal([]byte(v), &PrivateChat[k])
		if err != nil {
			fmt.Println("in func ProcessGetUnReadPvMsg,json.Unmarshal failed!", err)
			ResponseGetUnReadPvMsg(conn, msg.ID,
				model.UnReadMsgRes{Error: model.Err_InternalServerError.Error()})
			return
		}
	}
	resMsg.Data = PrivateChat
	err = ResponseGetUnReadPvMsg(conn, msg.ID, resMsg)
	if err != nil {
		fmt.Println("响应客户端失败！", err)
	}
	_,err=redisConn.Do("del", key)
	if err != nil {
		fmt.Println("数据库删除消息失败！",err)
	}
	fmt.Println("用户ID:",msg.ID,"请求查询未读私聊消息，","查询成功")
}

func ResponseGetUnReadPvMsg(conn net.Conn, userId uint64, resMsg model.UnReadMsgRes) error {
	_, ok := model.OnlineUsers[userId]
	if !ok {
		return model.Err_UserNotOnline
	}
	msg,err:=MessagePacking(model.UnReadMsgResType,resMsg)
	if err != nil {
		return err
	}
	err=SendMessage(conn,msg)
	if err != nil {
		return err
	}
	return nil
}
