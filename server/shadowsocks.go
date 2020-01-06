package main

import (
	"bufio"
	"encoding/binary"
	"github/wziww/medusa"
	"github/wziww/medusa/log"
	"net"
	"sync"
)

var (
	bufsize uint16 = 32 << 10
)

var bp sync.Pool

func init() {
	bp.New = func() interface{} {
		b := make([]byte, 269)
		return b
	}
}

func btsPoolGet() []byte {
	return bp.Get().([]byte)
}

func btsPoolPut(b []byte) {
	bp.Put(b)
}
func sshandleConn(conn *medusa.TCPConn) {
	defer func() {
		conn.Close()
	}()
	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 255(max length address) + 2(port) + 10(hmac-sha1)
	buf := btsPoolGet()
	defer func() {
		btsPoolPut(buf)
	}()
	n, _ := conn.Read(buf)
	ivlen := (*conn.Encryptor).Ivlen()
	if n < ivlen {
		log.FMTLog(log.LOGDEBUG, "package len error")
		return
	}
	buf = (*conn.Encryptor).Decode(buf[ivlen:n], buf[:ivlen])
	if len(buf) < (1 + 1 + 2) { // 1(addrType) + 1(lenByte) + 2(port)
		log.FMTLog(log.LOGDEBUG, "package len error")
		return
	}
	n = 0
	/**
	   The localConn connects to the dstServer, and sends a ver
	   identifier/method selection message:
		          +----+----------+----------+
		          |VER | NMETHODS | METHODS  |
		          +----+----------+----------+
		          | 1  |    1     | 1 to 255 |
		          +----+----------+----------+
	   The VER field is set to X'05' for this ver of the protocol.  The
	   NMETHODS field contains the number of method identifier octets that
	   appear in the METHODS field.
	*/
	// 第一个字段VER代表Socks的版本，Socks5默认为0x05，其固定长度为1个字节

	var dIP []byte
	// aType 代表请求的远程服务器地址类型，值长度1个字节，有三种类型
	switch buf[0] {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = buf[1 : 1+net.IPv4len]
		n = 1 + net.IPv4len
	case 0x03:
		//	DOMAINNAME: X'03'
		n = 2 + int(buf[1])
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[2:n]))
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		//	IP V6 address: X'04'
		dIP = buf[1 : 1+net.IPv6len]
		n = 1 + net.IPv6len
	default:
		return
	}
	dPort := buf[n : n+2]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	// 连接真正的远程服务
	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	}
	defer dstServer.Close()
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	dstServer.SetLinger(0)
	_, werr := dstServer.Write(buf[n+2:]) // 头部信息之外数据写入
	if werr != nil {
		log.FMTLog(log.LOGERROR, werr)
		return
	}
	// 响应客户端连接成功
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	// 响应客户端连接成功
	// 进行转发
	// 从 localUser 读取数据发送到 dstServer
	go func() {
		err := conn.SSDecodeCopy(&medusa.TCPConn{
			L:         dstServer.LocalAddr().String(),
			R:         dstServer.RemoteAddr().String(),
			Closer:    dstServer,
			Writer:    dstServer,
			Reader:    bufio.NewReader(dstServer),
			Encryptor: conn.Encryptor,
		}, nil) //
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			conn.Close()
			dstServer.Close()
		}
	}()
	// 从 dstServer 读取数据发送到 localUser，这里因为处在翻墙阶段出现网络错误的概率更大
	(&medusa.TCPConn{
		L:         dstServer.LocalAddr().String(),
		R:         dstServer.RemoteAddr().String(),
		Writer:    dstServer,
		Closer:    dstServer,
		Reader:    bufio.NewReader(dstServer),
		Encryptor: conn.Encryptor,
	}).SSEncodeCopy(conn)
}
