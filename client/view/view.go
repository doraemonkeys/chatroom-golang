package view

import (
	"bufio"
	"client/model"
	"fmt"
	"os"
	"strings"
)

func ShowMemu() {
	fmt.Println("---------欢迎使用多人聊天系统--------")
	fmt.Println("1.登录")
	fmt.Println("2.注册")
	fmt.Println("3.退出")
	fmt.Println("请选择<1-3>:")
}

func Home() {
	fmt.Println("1.显示在线用户列表")
	fmt.Println("2.发送群聊消息")
	fmt.Println("3.发送私聊消息")
	fmt.Println("4.查询本地私聊记录")
	fmt.Println("5.查询未读私聊消息")
	fmt.Println("6.查询群聊消息记录")
	fmt.Println("7.查询用户信息")
	fmt.Println("8.退出登录")
	fmt.Println("tips:聊天时输入\"EXIT\"可以回到主界面哦！")
	fmt.Println()
}

//用于在终端上展示用户信息
func ShowUserInfo(user model.UserInfo) error {
	time := strings.Split(user.RegisterTime, ".")
	fmt.Printf("-----------------------------------------\n")
	fmt.Printf("|  ID: %v                     \n", user.UserId)
	fmt.Printf("|  Name: %s    SEX: %s        \n", user.UserName, user.Sex)
	fmt.Printf("|  %s                         \n", time[0])
	fmt.Printf("|  Status:%s                  \n", user.Status)
	fmt.Printf("-----------------------------------------\n")
	return nil
}

//用于在终端上展示在线用户列表
func ShowOnlineList(onlineList []string) {
	fmt.Printf("----------------在线用户------------------\n")
	for _, v := range onlineList {
		fmt.Printf("| %s    \n", v)
	}
	fmt.Printf("-----------------------------------------\n")
}

func ShowGroupChatMsg(msg model.GroupChatMsg) {
	fmt.Println()
	fmt.Println("(群聊)", msg.SrcUser.UserName, "ID:", msg.SrcUser.UserId, msg.ChatTime)
	fmt.Printf("----->[%s]", msg.ChatData)
	fmt.Println()
}

func ShowPrivateChatMsg(msg model.PrivateChatMsg) {
	fmt.Println()
	fmt.Println("(私聊)", msg.SrcUser.UserName, "ID:", msg.SrcUser.UserId, msg.ChatTime)
	fmt.Printf("----->[%s]", msg.ChatData)
	fmt.Println()
}

func ShowPrivateChatMsgHistory(msg []model.PrivateChatMsg) {
	fmt.Println("按ENTER键浏览，输入EXIT退出")
	fmt.Println("------------------------------------------------")
	buf := bufio.NewReader(os.Stdin)
	for _, v := range msg {
		ShowPrivateChatMsg(v)
		input, _ := buf.ReadString('\n')
		input = input[:len(input)-2]
		if len(input) == 4 && strings.ToUpper(input) == "EXIT" {
			return
		}
	}
}

func ShowUnReadPvMsg(msg model.UnReadMsgRes) {
	if msg.Error != "" {
		fmt.Println("获取未读消息失败!", msg.Error)
		return
	}
	fmt.Println("按ENTER键浏览，输入EXIT退出")
	fmt.Println("------------------------------------------------")
	buf := bufio.NewReader(os.Stdin)
	for _, v := range msg.Data {
		ShowPrivateChatMsg(v)
		input, _ := buf.ReadString('\n')
		input = input[:len(input)-2]
		if len(input) == 4 && strings.ToUpper(input) == "EXIT" {
			return
		}
	}
}

func ShowUnReadMsgNotification(msg model.UnReadMsgNotification) {
	fmt.Println()
	fmt.Printf("-----------------------------------------\n")
	fmt.Printf("|****[通知] 您有%d 条未读私聊消息\n", msg.MsgNum)
	fmt.Printf("-----------------------------------------\n")
}
