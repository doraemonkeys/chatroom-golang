package server

import (
	"bufio"
	"client/model"
	"client/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func GetUserInfo(conn net.Conn) error {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("请输入要查询的ID:")
	var queryID string
	for {
		input, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		queryID = strings.TrimSpace(input)
		if !utils.UserIDTest(queryID) {
			fmt.Println("输入的ID有误，请重新输入！")
			continue
		}
		break
	}
	var ID uint64 = utils.Str2uint64(queryID)

	err := GetUserInfoByID(conn, ID)
	if err != nil {
		return err
	}
	return nil
}

//向服务端发起查询ID请求
func GetUserInfoByID(conn net.Conn, userID uint64) error {
	var GetUsreInfoMsg = model.GetUsreInfoMsg{
		UserId: userID,
	}
	data, err := json.Marshal(GetUsreInfoMsg)
	if err != nil {
		return err
	}
	var msg = model.Messsage{
		MsgType: model.GetUsreInfoType,
		Data:    string(data),
	}
	msg.DataLength = len(msg.Data)
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}
	return nil
}

//向服务端发起显示在线用户列表的(目前限制返回最多10个用户)
func GetOnlineList(conn net.Conn) error {
	var msg = model.Messsage{
		MsgType: model.GetOnlineUsersType,
		Data:    "",
	}
	msg.DataLength = len(msg.Data)
	err := SendMessage(conn, msg)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalPrivateChatHistory() ([]model.PrivateChatMsg, error) {
	fmt.Println("输入你要聊天记录查询的ID：")
	targetID := ""
	for {
		fmt.Scan(&targetID)
		if !utils.UserIDTest(targetID) {
			fmt.Println("输入的ID有误，请重新输入！")
			continue
		}
		break
	}
	fmt.Scanf("\n")
	name := fmt.Sprintf("CacheFiles/%d/%s.txt", LoggedInUser.UserId, targetID)
	path := fmt.Sprintf("CacheFiles/%d", LoggedInUser.UserId)
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0644)
		return nil, err
	}
	records, err := utils.ReverseRead(name, 100)
	if err != nil {
		return nil, err
	}
	recordMsgs := make([]model.PrivateChatMsg, 0, 110)
	var temMsg model.PrivateChatMsg
	for _, v := range records {
		if v == "" || v[0] == ' ' || v == " " { //跳过文件中可能的不规范行
			continue
		}
		err := json.Unmarshal([]byte(v), &temMsg)
		if err != nil {
			return nil, err
		}
		recordMsgs = append(recordMsgs, temMsg)
	}
	return recordMsgs, nil
}
