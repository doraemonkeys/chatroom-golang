package controller

import (
	"errors"
	"fmt"
	"net"
	"server/model"
	"server/server"
	"time"
)

//用户必须先登录,或者注册后登录
func LoginOrRegister(conn net.Conn) (uint64, error) {
	loginCount := 0 //客户端尝试登录次数(注册ip次数限制暂时没写)
	for {
		//接收客户端登录或注册消息
		msg, err := server.ReceiveMessage(conn)
		if err == model.Err_ConnectionDown {
			return 0, err //客户端断开连接直接返回
		}
		if err != nil {
			fmt.Println("ReceiveMessage failed,err:", err)
			continue
		}
		msgData, err := server.ParseMessage(msg)
		if err != nil {
			fmt.Println("in func LoginOrRegister,model.ParseMessage failed,err:", err)
			continue
		}
		switch msg.MsgType {
		//接收到登录消息
		case model.LoginMsgType:
			if loginCount > 5 {
				//TODO:限制对应ID 30s后再试
				server.ResponseLogin(conn, model.Err_FrequentlyLogin)
				err := server.RestrictLogin(msgData.(model.LoginMsg).UserId, conn.RemoteAddr().String())
				if err != nil {
					fmt.Println("in func LoginOrRegister,server.RestrictLogin failed")
				}
				loginCount = 0
				break
			}
			loginCount++
			err := server.Login(conn, msgData.(model.LoginMsg))
			if err != nil {
				fmt.Println("登录失败,err:", err)
				break
			}
			//登录成功
			return msgData.(model.LoginMsg).UserId, nil
		//接收到注册消息
		case model.RegisterMsgType:
			err := server.Register(conn, msgData.(model.RegisterMsg))
			if err != nil {
				fmt.Println("注册失败,err:", err)
			}
		default:
			fmt.Println("in func MainProcess,消息类型不正确")
			err = errors.New("消息类型不正确，请先登录！")
			server.ResponseLogin(conn, err)
		}
	}
}

//根据消息类型处理消息
func ProcessMessage(conn net.Conn, msg model.Messsage, UserId uint64) {
	msgData, err := server.ParseMessage(msg)
	if err != nil {
		fmt.Println("in func ProcessMessage,ParseMessage failed,err", err)
		return
	}
	switch msg.MsgType {
	case model.GetUsreInfoType:
		err := server.ProcessGetUsreInfo(conn, msgData.(model.GetUsreInfoMsg))
		if err != nil {
			fmt.Println("用户ID", UserId, "请求查询ID:", msgData.(model.GetUsreInfoMsg).UserId, "查询失败,err:", err)
			return
		}
		fmt.Println("用户ID", UserId, "请求查询ID:", msgData.(model.GetUsreInfoMsg).UserId, "查询成功")
	case model.GetOnlineUsersType:
		err := server.ProcessGetOnlineUsers(conn)
		if err != nil {
			fmt.Println("用户ID", UserId, "请求查询在线用户列表")
			fmt.Println("获取在线用户列表失败！err:", err)
		}
		fmt.Println("用户ID", UserId, "请求查询在线用户列表:", "查询成功")
	case model.GroupChatType:
		server.ProcessGroupChat(conn, msgData.(model.GroupChatMsg))
	case model.PrivateChatType:
		server.ProcessPrivateChat(conn, msgData.(model.PrivateChatMsg))
	case model.GetUnReadMsgType:
		server.ProcessGetUnReadPvMsg(conn, msgData.(model.GetUnReadMsg))
	case model.GetGroupChatType:
		server.ProcessGetGroupChatMsg(conn, msgData.(model.GetGroupChatMsg))
	default:
		fmt.Println("in func ProcessMessage,接收到不支持的消息类型,类型值:", msg.MsgType)
		return
	}
}

//用于维护用户在线列表，为保证map线程安全，通过对一个全局的channel发送指令来操作map。
//向通道CH_OnlineUser发送UpdateOperate指令可更新在线用户列表，然后由此函数处理。
//此函数并发不安全，只能在main中开一个goroutine
func UpdateOnlineUser() {
	for v := range model.CH_OnlineUser {
		switch v.OPcode {
		case model.Remove:
			delete(model.OnlineUsers, v.UserInfo.Info.UserId)
		case model.Add:
			model.OnlineUsers[v.UserInfo.Info.UserId] = v.UserInfo
		default:
			fmt.Println("in func UpdateOnlineUser,不支持的操作类型,code:", v.OPcode)
		}
	}
}

//总控
func MainProcess(conn net.Conn) {
	defer conn.Close()
	UserId, err := LoginOrRegister(conn)
	if err != nil {
		fmt.Println("ip:", conn.RemoteAddr().String(), "登录失败,err:", err)
		return
	}
	defer server.Logout(UserId)
	//到此用户登录成功
	fmt.Println()
	//登录成功后继续接收客户端的其他消息
	for {
		msg, err := server.ReceiveMessage(conn)
		if err == model.Err_ConnectionDown {
			return
		}
		if err != nil {
			fmt.Println("in func MainProcess,ReceiveMessage failed,err:", err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		go ProcessMessage(conn, msg, UserId)
	}
}
