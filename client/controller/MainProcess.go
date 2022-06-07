package controller

import (
	"bufio"
	"client/model"
	"client/server"
	"client/view"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)



func MainProcess(conn net.Conn, userInfo model.UserInfo) {
	view.Home() //显示主界面
	go RecMsg(conn, userInfo)
	buf := bufio.NewReader(os.Stdin)
	//处理用户的输入指令
	for {
		fmt.Printf("请选择<1-8>:")
		input, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("in func MainProcess,读取用户输入失败,err:", err)
		}
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			err := server.GetOnlineList(conn)
			if err != nil {
				fmt.Println("显示在线用户列表遇到问题，err", err)
			}
			ok := server.CheckResponse(model.OnlineUsersResType)
			if !ok {
				fmt.Println("服务端响应超时！")
			}
		case "2":
			server.SendGroupChat(conn, userInfo)
		case "3":
			server.SendPrivateChat(conn, userInfo)
		case "4":
			msg, err := server.GetLocalPrivateChatHistory()
			if err != nil {
				fmt.Println("查找失败!", err)
				break
			}
			view.ShowPrivateChatMsgHistory(msg)
		case "5":
			err := server.GetUnReadPrivateMsg(conn, userInfo.UserId)
			if err != nil {
				fmt.Println("获取未读消息失败!", err)
			}
		case "6":
			err := server.GetGroupChatMsg(conn, userInfo)
			if err != nil {
				fmt.Println("获取群聊消息失败!", err)
			}
		case "7":
			err := server.GetUserInfo(conn)
			if err != nil {
				fmt.Println("查询用户失败！err:", err)
			}
			ok := server.CheckResponse(model.GetUsreInfoResType)
			if !ok {
				fmt.Println("服务端响应超时！")
			}
		case "8":
			return
		case "exit":
			return
		case "EXIT":
			return
		default:
			fmt.Println("你的输入有误，请输入<1-8>:")
		}
	}
}

//循环接收服务器发出的消息
func RecMsg(conn net.Conn, userInfo model.UserInfo) {
	for {
		NewMsg, err := server.ReceiveMessage(conn)
		if err != nil {
			if err == io.EOF || err == model.Err_ConnectionDown {
				return
			}
			fmt.Println("in func ProcessMsg,ReceiveMessage failed,err:", err)
			continue
		}
		//可能会同时接收很多消息，所以开goroutine处理
		go ProcessMsg(NewMsg)
	}
}

//处理接收到的消息
func ProcessMsg(NewMsg model.Messsage) {
	server.MsgResType = NewMsg.MsgType
	MsgData, err := server.ParseMessage(NewMsg)
	if err != nil {
		fmt.Println("in func ProcessMsg,解析消息错误,err:", err)
		return
	}
	switch NewMsg.MsgType {
	case model.GetUsreInfoResType:
		if MsgData.(model.GetUsreInfoResMsg).Error != "" {
			fmt.Println("获取用户信息失败，err", MsgData.(model.GetUsreInfoResMsg).Error)
			return
		}
		err := view.ShowUserInfo(MsgData.(model.GetUsreInfoResMsg).UserInfos)
		if err != nil {
			fmt.Println("in func ProcessMsg,展示用户信息失败,err", err)
			return
		}
	case model.OnlineUsersResType:
		if MsgData.(model.OnlineUsersResMsg).Error != "" {
			fmt.Println("获取在线用户列表失败，err", MsgData.(model.OnlineUsersResMsg).Error)
			return
		}
		view.ShowOnlineList(MsgData.(model.OnlineUsersResMsg).Users)
	case model.GroupChatType:
		view.ShowGroupChatMsg(MsgData.(model.GroupChatMsg))
	case model.PrivateChatType:
		view.ShowPrivateChatMsg(MsgData.(model.PrivateChatMsg))
		err := server.SavePrivateChatMsg(MsgData.(model.PrivateChatMsg))
		if err != nil {
			fmt.Println("本地保存消息失败！", err)
		}
	case model.UnReadMsgNotificationType:
		view.ShowUnReadMsgNotification(MsgData.(model.UnReadMsgNotification))
	case model.UnReadMsgResType:
		model.Ch1<- MsgData.(model.UnReadMsgRes)
	case model.GetGroupChatResType:
		model.Ch1<- MsgData.(model.GetGroupChatResMsg)
	default:
		fmt.Println("接收到不支持的消息类型！")
		return
	}
}
