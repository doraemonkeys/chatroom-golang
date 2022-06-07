package server

import (
	"client/model"
	"client/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

//全局变量，当前登录的用户
var LoggedInUser model.UserInfo

func Login(conn net.Conn) (model.UserInfo, error) {
	LoginMsg, err := Get_LoginUser()
	if err != nil {
		fmt.Println("Get_LoginUser failed")
		return model.UserInfo{}, err
	}
	err = LoginUser(conn, LoginMsg)
	if err != nil {
		return model.UserInfo{}, err
	}
	fmt.Println("登录成功，请稍候...")
	userInfoMsg, err := ReceiveMessage(conn)
	if err != nil {
		return model.UserInfo{}, err
	}
	msgData, err := ParseMessage(userInfoMsg)
	if err != nil {
		return model.UserInfo{}, err
	}
	return msgData.(model.GetUsreInfoResMsg).UserInfos, nil
}

//与服务端进行登录验证
func LoginUser(conn net.Conn, user model.LoginMsg) error {
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("in func LoginUser,json.Marshal failed")
		return err
	}
	msg := model.Messsage{
		MsgType: model.LoginMsgType,
		Data:    string(data),
	}
	msg.DataLength = len(msg.Data)
	err = SendMessage(conn, msg)
	if err != nil {
		fmt.Println("in func LoginUser,SendMessage failed")
		return err
	}
	RecMsg, err := ReceiveMessage(conn)
	if err != nil {
		fmt.Println("in func LoginUser,ReceiveMessage failed")
		return err
	}
	msgData, err := ParseMessage(RecMsg)
	if err != nil {
		return err
	}
	if msgData.(model.LoginMsgRes).StatusCode == model.LoginSucess {
		return nil
	}
	return errors.New(msgData.(model.LoginMsgRes).Error)
}

//用户输入登录ID与密码
func Get_LoginUser() (model.LoginMsg, error) {
	var Id string
	var pwd string
	fmt.Printf("请输入账号(ID):\n")
	for {
		fmt.Scan(&Id)
		if !utils.UserIDTest(Id) {
			fmt.Println("输入的ID有误，请重新输入！")
			continue
		}
		break
	}
	fmt.Printf("请输入账号密码:\n")
	fmt.Scan(&pwd)
	userId := utils.Str2uint64(Id)
	fmt.Scanf("\n") //干掉回车
	return model.LoginMsg{UserId: userId, UserPwd: pwd}, nil
}

func InitProfile(user model.UserInfo) error {
	LoggedInUser = user
	name := fmt.Sprintf("CacheFiles/%d", user.UserId)
	//查看文件夹状态,没有则创建
	_, err := os.Stat(name)
	if err != nil {
		err := os.MkdirAll(name, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
