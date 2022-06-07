package server

import (
	"net"
)

//连接到服务器
func ConnectToServer() (net.Conn, error) {
	//return net.Dial("tcp", "192.168.1.104:7777")
	return net.Dial("tcp", "127.0.0.1:7777")
}
