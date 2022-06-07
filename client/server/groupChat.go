package server

import (
	"bufio"
	"client/model"
	"client/view"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func SendGroupChat(conn net.Conn, myself model.UserInfo) {
	buf := bufio.NewReader(os.Stdin)
	//处理用户的输入指令
	fmt.Printf("你想对大家说些什么？按ENTER发送，输入EXIT退出\n")
	for {
		input, err := buf.ReadString('\n')
		input = input[:len(input)-2]
		if len(input) == 4 && strings.ToUpper(input) == "EXIT" {
			return
		}
		//fmt.Printf("调试,|%s|\n", input)
		if err != nil {
			fmt.Println("发送失败！err:", err)
			continue
		}
		var chatMsg = model.GroupChatMsg{
			SrcUser:  model.SimpleUserInfo{UserId: myself.UserId, UserName: myself.UserName},
			ChatData: input,
		}
		msg, err := MessagePacking(model.GroupChatType, chatMsg)
		if err != nil {
			fmt.Println("发送失败！err:", err)
			continue
		}
		err = SendMessage(conn, msg)
		if err != nil {
			fmt.Println("发送失败！err:", err)
			continue
		}
	}
}

func GetGroupChatMsg(conn net.Conn, user model.UserInfo) error {
	var MsgData = model.GetGroupChatMsg{
		SrcUser: model.SimpleUserInfo{UserId: user.UserId, UserName: user.UserName},
	}
	msg, err := MessagePacking(model.GetGroupChatType, MsgData)
	if err != nil {
		return err
	}
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}

	ResMsg := <-model.Ch1
	Msgs, ok := ResMsg.(model.GetGroupChatResMsg)
	if !ok {
		fmt.Println("in func GetGroupChatMsg,从通道接收到未知类型消息！")
		return model.Err_GetUnknowMsg
	}
	if Msgs.Error != "" {
		return errors.New(Msgs.Error)
	}
	fmt.Println("按ENTER键浏览，输入EXIT退出")
	fmt.Println("------------------------------------------------")
	buf := bufio.NewReader(os.Stdin)
	for _, v := range Msgs.Data {
		view.ShowGroupChatMsg(v)
		input, _ := buf.ReadString('\n')
		input = input[:len(input)-2]
		if len(input) == 4 && strings.ToUpper(input) == "EXIT" {
			return nil
		}
	}
	return nil
}
