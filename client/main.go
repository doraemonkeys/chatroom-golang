package main

import (
	"bufio"
	"client/controller"
	"client/server"
	"client/view"
	"fmt"
	"os"
	"strings"
)

func main() {
	conn, err := server.ConnectToServer()
	if err != nil {
		fmt.Println("连接服务器失败,err:", err)
		return
	}
	//初始界面读取用户输入
	buf := bufio.NewReader(os.Stdin)
	for {
		view.ShowMemu()
		input, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("func MainProcess,ReadString failed,err:", err)
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			userinfo, err := server.Login(conn)
			if err != nil {
				fmt.Println("登录失败，err:", err)
				fmt.Println("-------------按ENTER键继续----------------")
				buf.ReadString('\n')
				break
			}
			err = server.InitProfile(userinfo)
			if err != nil {
				fmt.Println("客户端初始化失败，请检查客户端权限!", err)
				break
			}
			fmt.Println("----------------------------------------")
			fmt.Println("恭喜", userinfo.UserName, "登录成功！")
			//进入主界面
			controller.MainProcess(conn, userinfo)
			conn.Close()
			conn, err = server.ConnectToServer()
			if err != nil {
				fmt.Println("连接服务器失败,err:", err)
				return
			}
		case "2":
			err := server.Register(conn)
			if err != nil {
				fmt.Println("注册失败，err:", err)
				break
			}
			fmt.Println("注册成功！")
			fmt.Println("------------------------------------")
		case "3":
			conn.Close()
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新输入！")
		}
	}
}
