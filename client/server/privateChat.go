package server

import (
	"bufio"
	"client/model"
	"client/utils"
	"client/view"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func SendPrivateChat(conn net.Conn, myself model.UserInfo) {
	fmt.Println("请输入私聊对象的ID:")
	DstUser := ""
	for {
		fmt.Scan(&DstUser)
		if !utils.UserIDTest(DstUser) {
			fmt.Println("输入的ID有误，请重新输入！")
			continue
		}
		break
	}
	fmt.Scanf("\n")
	buf := bufio.NewReader(os.Stdin)
	//处理用户的输入指令
	fmt.Printf("按ENTER发送，输入EXIT退出\n")
	for {
		input, err := buf.ReadString('\n')
		//fmt.Printf("调试,|%s|\n", input)
		if err != nil {
			fmt.Println("发送失败！err:", err)
			continue
		}
		input = input[:len(input)-2]
		if len(input) == 4 && strings.ToUpper(input) == "EXIT" {
			return
		}
		var chatMsg = model.PrivateChatMsg{
			SrcUser:  model.SimpleUserInfo{UserId: myself.UserId, UserName: myself.UserName},
			DstUser:  utils.Str2uint64(DstUser),
			ChatTime: time.Now().Format("2006.01.02 15:04"),
			ChatData: input,
		}
		err = SaveMyPvChat(chatMsg)
		if err != nil {
			fmt.Println("本地保存聊天失败！", err)
		}
		msg, err := MessagePacking(model.PrivateChatType, chatMsg)
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

//保存接收到的聊天记录
func SavePrivateChatMsg(msg model.PrivateChatMsg) error {
	name := fmt.Sprintf("CacheFiles/%d/%d.txt", msg.DstUser, msg.SrcUser.UserId)
	path := fmt.Sprintf("CacheFiles/%d", msg.DstUser)
	//查看文件夹状态,没有则创建
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0644)
	}
	//打开文件,没有则创建
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	file.Write(data)
	_, err = file.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

//保存我发送的聊天记录
func SaveMyPvChat(msg model.PrivateChatMsg) error {
	name := fmt.Sprintf("CacheFiles/%d/%d.txt", msg.SrcUser.UserId, msg.DstUser)
	path := fmt.Sprintf("CacheFiles/%d", msg.SrcUser.UserId)
	//查看文件夹状态,没有则创建
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0644)
	}
	//打开文件,没有则创建
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	file.Write(data)
	_, err = file.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

func GetUnReadPrivateMsg(conn net.Conn, userID uint64) error {
	var msgData = model.GetUnReadMsg{
		ID: userID,
	}
	msg, err := MessagePacking(model.GetUnReadMsgType, msgData)
	if err != nil {
		return err
	}
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}

	ResMsg := <-model.Ch1
	unRead, ok := ResMsg.(model.UnReadMsgRes)
	if !ok {
		fmt.Println("in func GetUnReadPrivateMsg,从通道接收到未知类型消息！")
		return model.Err_GetUnknowMsg
	}
	if unRead.Error!="" {
		return errors.New(unRead.Error)
	}
	for _, v := range unRead.Data {
		err := SavePrivateChatMsg(v)
		if err != nil {
			return err
		}
	}
	view.ShowUnReadPvMsg(unRead)
	return nil
}
