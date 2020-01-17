package main

import (
	"bufio"
	"github/wziww/medusa"
	"github/wziww/medusa/config"
	"github/wziww/medusa/encrypt"
	"github/wziww/medusa/log"
	"github/wziww/medusa/stream"
	"net"
	"os"
	"strconv"
)

func main() {
	// 配置初始化
	config.Init()
	// 日志初始化
	log.Init()
	addr, resoveErr := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(config.C.Server.Port))
	if resoveErr != nil {
		log.FMTLog(log.LOGERROR, resoveErr)
		os.Exit(0)
	}
	config.C.Base.Client = false
	log.FMTLog(log.LOGINFO, "server start")
	// api 服务初始化
	stream.APIServerInit()
	// 加密器初始化
	password := []byte(config.C.Base.Password)
	encryptor := encrypt.InitEncrypto(&password, config.C.Base.Crypto, config.C.Base.Padding, nil)
	if encryptor == nil {
		log.FMTLog(log.LOGERROR, "encrypto:", config.C.Base.Crypto, " init error,please checkout your config file")
		os.Exit(0)
	}
	// 服务启动
	listener, listenErr := net.ListenTCP("tcp", addr)
	if listenErr != nil {
		log.FMTLog(log.LOGERROR, listenErr)
		os.Exit(0)
	}
	log.FMTLog(log.LOGINFO, "service start listen at", addr)
	defer listener.Close()
	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			continue
		}
		if !config.CheckIPAllow(localConn.RemoteAddr().String()) {
			if config.CheckIPNotAllow(localConn.RemoteAddr().String()) {
				log.FMTLog(log.LOGINFO, localConn.RemoteAddr().String(), "connect refused")
				localConn.Close()
				continue
			}
		}
		// log.FMTLog(log.LOGINFO, localConn.RemoteAddr(), "connected")
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		encryptor := encrypt.InitEncrypto(&password, config.C.Base.Crypto, config.C.Base.Padding, []byte(encrypt.GetRandString(encryptor.Ivlen())))
		if encryptor == nil {
			localConn.Close()
			log.FMTLog(log.LOGERROR, "encrypto:", config.C.Base.Crypto, " init error,please checkout your config file")
			continue
		}
		go sshandleConn(&medusa.TCPConn{
			L:         localConn.LocalAddr().String(),
			R:         localConn.RemoteAddr().String(),
			Reader:    bufio.NewReader(localConn),
			Closer:    localConn,
			Writer:    localConn,
			Encryptor: &encryptor,
		})
	}
}
