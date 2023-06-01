package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// 用户名
var loginName string

// 本机连接
var selfConnect *net.TCPConn

// 读取行文本
var reader = bufio.NewReader(os.Stdin)

var quitCh = make(chan bool)

// 建立连接
func connect(addr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr) // 使用tcp
	con, err := net.DialTCP("tcp", nil, tcpAddr)  // 拨号：主动向server建立连接
	selfConnect = con
	if err != nil {
		fmt.Println("连接服务器失败")
		os.Exit(1)
	}
	go msgSender()
	go msgReceiver()
}

// 消息接收器
func msgReceiver() {
	defer func() {
		quitCh <- true
	}()
	buff := make([]byte, 2048)
	for {
		len, err := selfConnect.Read(buff) // 从建立连接的缓冲区读消息
		if err != nil {
			fmt.Printf("recv error: %s\n", err)
			return
		}

		if len == 0 {
			fmt.Printf("recv 0 bytes from server, server close the connection, quit the goroutine\n")
			return
		}

		fmt.Printf("recv %v bytes from server\n", len)
		fmt.Println(string(buff[:len]))
	}
}

// 消息发送器
func msgSender() {
	defer func() {
		quitCh <- true
	}()

	for {
		bMsg, _, _ := reader.ReadLine()
		bMsg = []byte(loginName + " : " + string(bMsg))
		_, err := selfConnect.Write(bMsg) // 发消息
		if err != nil {
			return
		}
	}
}

// 初始化
func initGroupChatClient() {
	fmt.Println("请问您怎么称呼？")
	bName, _, _ := reader.ReadLine()
	loginName = string(bName)
	connect("127.0.0.1:1801")
	<-quitCh
}

func main() {
	defer func() {
		if selfConnect != nil {
			selfConnect.Close()
		}

		fmt.Println("client quitted")
	}()
	initGroupChatClient()

}
