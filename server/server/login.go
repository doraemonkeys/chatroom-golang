package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"

	"github.com/garyburd/redigo/redis"
)

func LoginUser(msg model.LoginMsg) error {
	redisConn := redisPool.Get()
	data, err := redis.String(redisConn.Do("hget", "users", msg.UserId))
	if err != nil {
		fmt.Println("in func LoginUser,找不到用户", err)
		return model.Err_UserNotExist
	}
	var userData model.User
	json.Unmarshal([]byte(data), &userData)
	if utils.StrMatch(userData.UserPwd, msg.UserPwd) {
		return nil
	}
	return model.Err_PasswordWrong
}

//登录结果返回
func ResponseLogin(conn net.Conn, err error) error {
	var NewResMsg model.LoginMsgRes
	if err != nil {
		NewResMsg.StatusCode = model.LoginFailed
		NewResMsg.Error = err.Error()
	} else {
		NewResMsg.StatusCode = model.LoginSucess
		NewResMsg.Error = ""
	}
	data, err := json.Marshal(NewResMsg)
	if err != nil {
		fmt.Println("in func ResponseLogin,json.Marshal failed")
		return err
	}
	msg := model.Messsage{
		MsgType: model.LoginMsgResType,
		Data:    string(data),
	}
	msg.DataLength = len(msg.Data)
	err = SendMessage(conn, msg)
	if err != nil {
		fmt.Println("in func ResponseLogin,SendMessage failed")
		return err
	}
	return nil
}

//尝试登录并返回登录状态给客户端
func Login(conn net.Conn, NewLoginMsg model.LoginMsg) error {
	redisConn := redisPool.Get()
	//查询ID是否在登录频繁的黑名单中
	idHave, err := redis.Bool(redisConn.Do("exists", NewLoginMsg.UserId))
	if err != nil {
		fmt.Println("in func Login,在redis中查询id失败")
		ResponseLogin(conn, model.Err_InternalServerError)
		return err
	}
	if idHave {
		ResponseLogin(conn, model.Err_FrequentlyLogin)
		return model.Err_FrequentlyLogin
	}
	//登录
	err = LoginUser(NewLoginMsg)
	if err != nil {
		fmt.Println("用户:", NewLoginMsg.UserId, "密码:", NewLoginMsg.UserPwd, "登录失败")
		ResponseLogin(conn, err)
		return err
	}
	userInfo, err := GetUserInfoByID(NewLoginMsg.UserId)
	if err != nil {
		fmt.Println("用户:", NewLoginMsg.UserId, "密码:", NewLoginMsg.UserPwd, "登录失败")
		ResponseLogin(conn, err)
		return err
	}
	//若用户已经登录
	if userInfo.Status == model.Online {
		fmt.Println("用户:", NewLoginMsg.UserId, "密码:", NewLoginMsg.UserPwd, "登录失败")
		ResponseLogin(conn, model.Err_RepeatLogin)
		return model.Err_RepeatLogin
	}
	var OP = model.UpdateOperate{
		OPcode: model.Add,
		UserInfo: model.OnlineUserInfo{
			Conn: conn,
			Info: userInfo,
		},
	}
	//更新在线列表
	model.CH_OnlineUser <- OP
	//响应客户端
	err = ResponseLogin(conn, err)
	if err != nil {
		err = ResponseLogin(conn, nil) //再次尝试回复
		fmt.Println("in func Login,回复客户端失败")
		return err
	}
	//返回登录的用户的个人信息
	err = ProcessGetUsreInfo(conn, model.GetUsreInfoMsg{UserId: NewLoginMsg.UserId})
	if err != nil {
		return err
	}
	fmt.Println("ID:", NewLoginMsg.UserId, "密码:", NewLoginMsg.UserPwd, "登录成功")
	//检查未读消息，并提醒客户端
	err = CheckUnResdMsg(conn, NewLoginMsg)
	if err != nil {
		fmt.Println("检查未读消息失败，请排查服务端故障！", err)
	}
	return nil
}

//限制用户30秒后才能登录
func RestrictLogin(userId uint64, ip string) error {
	redisConn := redisPool.Get()
	_, err := redisConn.Do("set", userId, ip)
	if err != nil {
		fmt.Println("in func RestrictLogin,写入redis失败")
		return err
	}
	_, err = redisConn.Do("expire", userId, 30)
	if err != nil {
		fmt.Println("in func RestrictLogin,设置延时删除失败")
		return err
	}
	return nil
}

func Logout(userID uint64) {
	model.CH_OnlineUser <- model.UpdateOperate{
		OPcode: model.Remove,
		UserInfo: model.OnlineUserInfo{
			Info: model.UserInfo{
				UserId: userID,
			},
		},
	}
}
