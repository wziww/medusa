package main

import (
	"github/wziww/medusa"
	"github/wziww/medusa/config"
	"github/wziww/medusa/log"
	"net"
	"os"
)

func handleConn(userConn *medusa.TCPConn) {
	defer userConn.Close()

	dstAddr, resoveErr := net.ResolveTCPAddr("tcp", config.C.Client.RemoteAddress)
	if resoveErr != nil {
		log.FMTLog(log.LOGERROR, resoveErr)
		os.Exit(0)
	}

	// 连接真正的远程服务
	proxyServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return
	}
	defer proxyServer.Close()
	proxyServerTCP := &medusa.TCPConn{
		L:               proxyServer.LocalAddr().String(),
		R:               proxyServer.RemoteAddr().String(),
		ReadWriteCloser: proxyServer,
		Encryptor:       userConn.Encryptor,
	}
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	proxyServer.SetLinger(0)

	// 进行转发
	// 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := proxyServerTCP.DecodeCopy(userConn)
		if err != nil {
			log.FMTLog(log.LOGDEBUG, err)
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	userConn.EncodeCopy(proxyServerTCP)
}
