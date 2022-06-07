package server

import (
	"client/model"
	"client/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func Register(conn net.Conn) error {
	NewUser, err := Get_RegisterUser()
	if err != nil {
		return err
	}
	err = RegisterUser(conn, NewUser)
	if err != nil {
		return err
	}
	return nil
}

//未校验ID与密码,服务端的用户已存在的校验也没有写
func Get_RegisterUser() (model.User, error) {
	var (
		userName  string
		ID        string
		userID    uint64
		userPwd   string
		sex       string
		selectNum string
		NewUser   model.User
	)
	fmt.Println("请输入一个用户昵称:")
	fmt.Scan(&userName)
	for !utils.UserNameTest(userName) {
		fmt.Println("输入的用户名有误，请重新输入！")
		fmt.Scan(&userName)
	}
	fmt.Println("请输入一个用户ID(纯数字):")
	fmt.Scan(&ID)
	for !utils.UserIDTest(ID) {
		fmt.Println("输入的用户ID有误，请重新输入！")
		fmt.Scan(&ID)
	}
	fmt.Println("请输入一个用户密码：")
	fmt.Scan(&userPwd)
	for !utils.PasswordTest(userPwd) {
		fmt.Println("输入的密码不规范，请重新输入！")
		fmt.Scan(&userPwd)
	}
	fmt.Println("请选择性别<1-3>：1.男 2.女 3.未知 ")
	fmt.Scan(&selectNum)
	ok := false
	for !ok {
		ok = true
		switch selectNum {
		case "1":
			sex = "男"
		case "2":
			sex = "女"
		case "3":
			sex = "未知"
		default:
			fmt.Println("你的输入有误，请输入<1-3>:")
			fmt.Scan(&selectNum)
			ok = false
		}
	}
	userID = utils.Str2uint64(ID)
	NewUser.UserId = userID
	NewUser.UserName = userName
	NewUser.UserPwd = userPwd
	NewUser.Sex = sex
	fmt.Scanf("\n")
	return NewUser, nil
}

//与服务端交换信息进行注册
func RegisterUser(conn net.Conn, NewUser model.User) error {
	var registerMsg = model.RegisterMsg{
		User: NewUser,
	}
	data, err := json.Marshal(registerMsg)
	if err != nil {
		return err
	}
	msg := model.Messsage{
		MsgType: model.RegisterMsgType,
		Data:    string(data),
	}
	msg.DataLength = len(msg.Data)
	err = SendMessage(conn, msg)
	if err != nil {
		return err
	}
	RecMsg, err := ReceiveMessage(conn)
	if err != nil {
		return err
	}
	msgData, err := ParseMessage(RecMsg)
	if err != nil {
		return err
	}
	//注册成功
	if msgData.(model.RegisterMsgRes).StatusCode == model.RegisterSucess {
		return nil
	}
	return errors.New(msgData.(model.RegisterMsgRes).Error)
}
