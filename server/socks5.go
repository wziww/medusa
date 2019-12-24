package main

import (
	"bufio"
	"encoding/binary"
	"github/wziww/medusa"
	"github/wziww/medusa/log"
	"net"
)

// https://www.ietf.org/rfc/rfc1928.txt
func handleConn(conn *medusa.TCPConn) {
	defer func() {
		conn.Close()
	}()

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
	_, buf, err := conn.DecodeRead()
	// 只支持版本5
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return
	}
	if len(buf) < 1 || buf[0] != 0x05 {
		return
	}

	/**
	   The dstServer selects from one of the methods given in METHODS, and
	   sends a METHOD selection message:

		          +----+--------+
		          |VER | METHOD |
		          +----+--------+
		          | 1  |   1    |
		          +----+--------+
	*/
	// 不需要验证，直接验证通过
	conn.EncodeWrite([]byte{0x05, 0x00})

	// /**
	//   +----+-----+-------+------+----------+----------+
	//   |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	//   +----+-----+-------+------+----------+----------+
	//   | 1  |  1  | X'00' |  1   | Variable |    2     |
	//   +----+-----+-------+------+----------+----------+
	// */

	n, buf, err := conn.DecodeRead()
	if err != nil || n < 7 {
		log.FMTLog(log.LOGERROR, err)
		return
	}

	if buf[1] != 0x01 {
		// 目前只支持 CONNECT
		return
	}

	var dIP []byte
	// aType 代表请求的远程服务器地址类型，值长度1个字节，有三种类型
	switch buf[3] {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		//	DOMAINNAME: X'03'
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		//	IP V6 address: X'04'
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
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

	// 响应客户端连接成功
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	// 响应客户端连接成功
	conn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	// 进行转发
	// 从 localUser 读取数据发送到 dstServer
	go func() {
		err := conn.DecodeCopy(&medusa.TCPConn{
			L:         dstServer.LocalAddr().String(),
			R:         dstServer.RemoteAddr().String(),
			Closer:    dstServer,
			Writer:    dstServer,
			Reader:    bufio.NewReader(dstServer),
			Encryptor: conn.Encryptor,
		}) //
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
	}).EncodeCopy(conn)
}
