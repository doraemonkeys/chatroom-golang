package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"
	"time"
)

//尝试注册并返回注册状态给客户端
func Register(conn net.Conn, NewRegisterMsg model.RegisterMsg) error {
	//注册用户
	err := RegisterUser(NewRegisterMsg.User)
	//返回注册状态给客户端
	ResponseRegister(conn, NewRegisterMsg.User.UserId, err)
	if err != nil {
		fmt.Println("ID", NewRegisterMsg.User.UserId, "注册失败")
		return err
	}
	fmt.Println("ID:", NewRegisterMsg.User.UserId, "注册成功")
	return nil
}

//注册账号
func RegisterUser(user model.User) error {
	fmt.Println("用户:", user.UserName, "ID", user.UserId, "密码:", user.UserPwd, "请求注册")
	_,err:=GetUserByID(user.UserId)
	if err!=model.Err_UserNotExist {
		fmt.Println("用户已存在，拒绝注册！")
		return model.Err_UserAlreadyExist
	}
	user.UserPwd = utils.StrEncrypt(user.UserPwd)
	user.RegisterTime = time.Now().String()
	redisConn := redisPool.Get()
	defer redisConn.Close()
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("in func RegisterUser,json.Marshal failed")
		return err
	}
	_, err = redisConn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("in func RegisterUser,redisConn.Do failed,err:")
		return err
	}
	return nil
}

//返回注册状态给客户端
//userId:请求注册的用户ID
//err: 注册返回的err,nil表示注册成功
func ResponseRegister(conn net.Conn, userId uint64, err error) error {
	NewResMsg := model.RegisterMsgRes{
		UserId: userId,
	}
	if err != nil {
		NewResMsg.StatusCode = model.RegisterFailed
		NewResMsg.Error = err.Error()
	} else {
		NewResMsg.StatusCode = model.RegisterSucess
		NewResMsg.Error = ""
	}
	data, err := json.Marshal(NewResMsg)
	if err != nil {
		fmt.Println("in func ResponseRegister,json.Marshal failed")
		return err
	}
	msg := model.Messsage{
		MsgType: model.RegisterMsgResType,
		Data:    string(data),
	}
	msg.DataLength = len(msg.Data)
	err = SendMessage(conn, msg)
	if err != nil {
		fmt.Println("in func ResponseRegister,SendMessage failed")
		return err
	}
	return nil
}
