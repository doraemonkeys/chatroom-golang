package main

import (
	"fmt"
	"net"
	"server/controller"
	"server/server"
)

func main() {
	fmt.Println("初始化Redis服务....")
	err := server.InitRedisService()
	if err != nil {
		fmt.Println("初始化Redis服务失败,err:", err)
		return
	}
	fmt.Println("服务器在7777端口监听....")
	//listen, err := net.Listen("tcp", "192.168.1.104:7777")
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("func MainProcess,Listen failed,err:", err)
	}
	fmt.Println("等待客户端来连接....")
	//用于维护用户在线列表
	go controller.UpdateOnlineUser()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("func MainProcess,listen.Accept failed,err:", err)
			continue
		}
		fmt.Println("连接成功，客户端:", conn.RemoteAddr().String())
		//建立一个协程处理连接
		go controller.MainProcess(conn)
	}
}
