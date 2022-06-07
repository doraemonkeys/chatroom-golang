package server

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"

	"github.com/garyburd/redigo/redis"
)

//仅用于服务端通过ID查找用户
func GetUserByID(userID uint64) (model.User, error) {
	redisConn := redisPool.Get()
	data, err := redis.String(redisConn.Do("hget", "users", userID))
	if err != nil {
		fmt.Println("找不到ID为:", userID, "的用户：", err)
		return model.User{}, model.Err_UserNotExist
	}
	var userData model.User
	err = json.Unmarshal([]byte(data), &userData)
	if err != nil {
		return userData, err
	}
	return userData, err
}

//仅用于服务端通过ID查找用户获得UserInfo
func GetUserInfoByID(userID uint64) (model.UserInfo, error) {
	redisConn := redisPool.Get()
	data, err := redis.String(redisConn.Do("hget", "users", userID))
	if err != nil {
		fmt.Println("找不到ID为:", userID, "的用户：", err)
		return model.UserInfo{}, model.Err_UserNotExist
	}
	var userData model.User
	err = json.Unmarshal([]byte(data), &userData)
	if err != nil {
		return model.UserInfo{}, err
	}
	var userInfo = model.UserInfo{
		UserId:       userData.UserId,
		UserName:     userData.UserName,
		Sex:          userData.Sex,
		RegisterTime: userData.RegisterTime,
		Status:       CheckUserStatus(userData.UserId),
	}
	return userInfo, err
}

//处理获取客户端用户信息的消息
func ProcessGetUsreInfo(conn net.Conn, msg model.GetUsreInfoMsg) error {
	userInfos, err := GetUserInfoByID(msg.UserId)
	if err != nil {
		NewErr := ResponseGetUsreInfo(conn, model.UserInfo{}, err)
		if NewErr != nil {
			fmt.Println("响应客户端失败，err", NewErr)
		}
		return err
	}
	err = ResponseGetUsreInfo(conn, userInfos, nil)
	if err != nil {
		return err
	}
	return nil
}

//回复获取的用户信息
func ResponseGetUsreInfo(conn net.Conn, userInfos model.UserInfo, err error) error {
	ResInfo := model.GetUsreInfoResMsg{
		UserInfos: userInfos,
		Error:     "",
	}
	if err != nil {
		ResInfo.Error = err.Error()
	}
	data, err := json.Marshal(ResInfo)
	if err != nil {
		return err
	}
	ResMsg := model.Messsage{
		MsgType: model.GetUsreInfoResType,
		Data:    string(data),
	}
	ResMsg.DataLength = len(ResMsg.Data)
	err = SendMessage(conn, ResMsg)
	if err != nil {
		return err
	}
	return nil
}

//检查用户状态，若用户的tcp连接不正常则判断为离线。
//若用户离线且在线列表中仍然存在其ID，则更新在线列表
func CheckUserStatus(userID uint64) model.UserStatus {
	v, exist := model.OnlineUsers[userID]
	if !exist {
		return model.Offline
	}
	_, err := v.Conn.Write(nil)
	if err != nil {
		Logout(userID)
		return model.Offline
	}
	return model.Online
}

func ProcessGetOnlineUsers(conn net.Conn) error {
	var onlineList = make([]string, 0, 10)
	count := 0
	for _, v := range model.OnlineUsers {
		ID := "ID: " + fmt.Sprintf("%v", v.Info.UserId)
		for len(ID) < 20 { //为了客户端排版整齐
			ID = ID + " "
		}
		Name := "Name: " + v.Info.UserName
		onlineList = append(onlineList, ID+Name)
		count++
		//限制显示10条
		if count == 10 {
			break
		}
	}
	if len(onlineList) == 0 {
		fmt.Println("逻辑上出现问题，怎么可能没有在线用户！")
		onlineList = append(onlineList, "无用户在线")
	}
	resList := model.OnlineUsersResMsg{
		Users: onlineList,
		Error: "",
	}
	if onlineList[0] == "无用户在线" {
		resList.Error = model.Err_NoOnlineUsers.Error()
	}
	data, err := json.Marshal(resList)
	if err != nil {
		return err
	}
	var ResMsg = model.Messsage{
		MsgType: model.OnlineUsersResType,
		Data:    string(data),
	}
	ResMsg.DataLength = len(ResMsg.Data)
	err = SendMessage(conn, ResMsg)
	if err != nil {
		return err
	}
	return nil
}
