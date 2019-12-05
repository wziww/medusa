package medusa

import (
	"encoding/binary"
	"errors"
	"github/wziww/medusa/log"
	"io"
)

var bufSize int = 1024

// TCPConn ...
type TCPConn struct {
	io.ReadWriteCloser
	Encryptor Encryptor
}

// DecodeRead ...
func (conn *TCPConn) DecodeRead() (n int, buf []byte, err error) {
	var l int64
	binary.Read(conn, binary.BigEndian, &l)
	data := make([]byte, l)
	n, err = conn.Read(data)
	if err != nil {
		log.FMTLog(log.LOGDEBUG, err)
		return
	}
	if res := conn.Encryptor.Decode(data[:n]); res != nil {
		buf = res
		n = len(res)
		return
	}
	return 0, nil, errors.New("DecodeRead error")
}

//EncodeWrite ...
func (conn *TCPConn) EncodeWrite(buf []byte) (n int, err error) {
	buf = conn.Encryptor.Encode(buf)
	if buf != nil {
		var l int64 = int64(len(buf))
		binary.Write(conn, binary.BigEndian, l)
		return conn.Write(buf)
	}
	return
}

// DecodeCopy 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) DecodeCopy(dst io.Writer) error {
	for {
		readCount, buf, errRead := conn.DecodeRead()
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// EncodeCopy 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := conn.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
		if readCount > 0 {
			writeCount, errWrite := (&TCPConn{
				ReadWriteCloser: dst,
				Encryptor:       conn.Encryptor,
			}).EncodeWrite(buf[:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount && writeCount == 0 {
				return io.ErrShortWrite
			}
		}
	}
}
