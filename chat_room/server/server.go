package main

import (
	"fmt"
	"net"
	"time"
)

// 客户端 map
var clientMap = make(map[string]*net.TCPConn) // 存储当前群聊中所有用户连接信息：key: ip+port, val: 用户连接信息

// 监听请求
func listenClient(ipAndPort string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ipAndPort)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for { // 循环接收
		clientConn, err := tcpListener.AcceptTCP() // 监听请求连接
		if err == nil {
			clientMap[clientConn.RemoteAddr().String()] = clientConn // 将连接添加到 map
		} else {
			fmt.Printf("Accept TCP faild: %s", err)
			continue
		}
		go addReceiver(clientConn)
		fmt.Println("用户 : ", clientConn.RemoteAddr().String(), " 已连接.")
	}
}

// 向连接添加接收器
func addReceiver(newConnect *net.TCPConn) {
	defer func() {
		delete(clientMap, newConnect.RemoteAddr().String())
		newConnect.Close()
	}()

	for {
		byteMsg := make([]byte, 2048)
		len, err := newConnect.Read(byteMsg) // 从newConnect中读取信息到缓存中

		if err != nil {
			fmt.Printf("recv error: %s\n", err)
			return
		}

		if len == 0 {
			fmt.Printf("recv 0 bytes from %s, client close the connection, quit the goroutine\n", newConnect.RemoteAddr().String())
			return
		}

		fmt.Printf("recv %v bytes from %s\n", len, newConnect.RemoteAddr().String())
		fmt.Println(string(byteMsg[:len]))
		msgBroadcast(byteMsg[:len], newConnect.RemoteAddr().String())
	}
}

// 广播给所有 client
func msgBroadcast(byteMsg []byte, key string) {
	for k, con := range clientMap {
		if k != key { // 转发消息给当前群聊中，除自身以外的其他用户
			con.Write(byteMsg)
		}
	}
}

// 初始化
func initGroupChatServer() {
	fmt.Println("服务已启动...")
	time.Sleep(1 * time.Second)
	fmt.Println("等待客户端请求连接...")
	go listenClient("127.0.0.1:1801")
	select {}
}

func main() {
	initGroupChatServer()
}
