package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"time"

	"github.com/garyburd/redigo/redis"
)

func ProcessGroupChat(conn net.Conn, msgData model.GroupChatMsg) {
	msgData.ChatTime = time.Now().Format("2006.01.02 15:04")
	data, err := json.Marshal(msgData)
	if err != nil {
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送群聊消息失败")
		return
	}
	redisConn := redisPool.Get()
	defer redisConn.Close()
	//目前没写删除群聊消息记录的逻辑
	redisConn.Do("rpush", "GroupChatMessage", string(data))
	var chatMsg = model.GroupChatMsg{
		SrcUser:  msgData.SrcUser,
		ChatData: msgData.ChatData,
		ChatTime: msgData.ChatTime,
	}
	msg, err := MessagePacking(model.GroupChatType, chatMsg)
	if err != nil {
		fmt.Println("用户ID:", msgData.SrcUser.UserId, "请求发送群聊消息失败")
		return
	}
	//给在线用户转发群聊消息
	for _, v := range model.OnlineUsers {
		err := SendMessage(v.Conn, msg)
		if err != nil {
			fmt.Println("转发群聊消息给", v.Info.UserId, "失败!", err)
		}
	}
	fmt.Println("用户ID:", msgData.SrcUser.UserId, "发送了一条群聊消息")
	fmt.Printf("内容为:[%s]\n", msgData.ChatData)
}

func ProcessGetGroupChatMsg(conn net.Conn, msg model.GetGroupChatMsg) {
	fmt.Println("用户ID:", msg.SrcUser.UserId, "请求查询群聊消息记录")
	redisConn := redisPool.Get()
	defer redisConn.Close()
	QueryNum := -10 //查询10条记录
	//目前没写删除群聊消息记录的逻辑
	data, err := redis.Strings(redisConn.Do("lrange", "GroupChatMessage", QueryNum, -1))
	if err != nil {
		fmt.Println("从数据库查询群聊消息失败！", err)
		ResponseGetGroupChatMsg(conn, model.GetGroupChatResMsg{Error: model.Err_InternalServerError.Error()})
		return
	}
	GroupMsgs := make([]model.GroupChatMsg, len(data))
	for k, v := range data {
		err := json.Unmarshal([]byte(v), &GroupMsgs[k])
		if err != nil {
			fmt.Println("in func ProcessGetGroupChatMsg,json.Unmarshal failed!", err)
			ResponseGetGroupChatMsg(conn, model.GetGroupChatResMsg{Error: model.Err_InternalServerError.Error()})
			return
		}
	}
	resMsg := model.GetGroupChatResMsg{
		Data:  GroupMsgs,
		Error: "",
	}
	err = ResponseGetGroupChatMsg(conn, resMsg)
	if err != nil {
		fmt.Println("响应客户端失败！", err)
	}
	fmt.Println("用户ID:", msg.SrcUser.UserId, "查询群聊消息记录，", "查询成功")
}

func ResponseGetGroupChatMsg(conn net.Conn, resMsg model.GetGroupChatResMsg) error {
	msg, err := MessagePacking(model.GetGroupChatResType, resMsg)
	if err != nil {
		return err
	}
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}
	return nil
}
