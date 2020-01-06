package medusa

import (
	"bufio"
	"github/wziww/medusa/encrypt"
	"github/wziww/medusa/stream"
	"io"
)

const tag = 0x80

var bufSize int64 = 32 << 10 //1kb

const maxConsecutiveEmptyReads = 100

// TCPConn ...
type TCPConn struct {
	L string // LocalAddr
	R string
	*bufio.Reader
	io.Closer
	io.Writer
	Encryptor *encrypt.Encryptor
	ivSent    *bool
}

// SSEncodeCopy 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) SSEncodeCopy(dst *TCPConn) error {
	buf := make([]byte, bufSize)
	for {
		c, err := ssEncodeCopy(conn, dst, buf)
		if !c {
			return err
		}
	}
}

// SSDecodeCopy 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) SSDecodeCopy(dst *TCPConn, iv []byte) error {
	buf := make([]byte, bufSize)
	for {
		c, err := ssDecodeCopy(conn, dst, buf, iv)
		if !c {
			return err
		}
	}
}
func (conn *TCPConn) ssEncodeWrite(buf []byte) (int, error) {
	buf = (*conn.Encryptor).Encode(buf)
	if conn.ivSent == nil {
		conn.ivSent = new(bool)
	}
	if !(*conn.ivSent) {
		*conn.ivSent = true
		buf = append((*conn.Encryptor).GetIv(), buf...)
	} else {
		b := make([]byte, len(buf))
		copy(b, buf)
	}
	var n int
	var e error
	for i := 0; i < maxConsecutiveEmptyReads; i++ {
		if n == len(buf) {
			break
		}
		nn, e := conn.Write(buf[n:])
		n += nn
		if e != nil {
			return n, e
		}
	}
	if n != len(buf) {
		return 0, io.ErrNoProgress
	}
	return n, e
}

func ssEncodeCopy(conn, dst *TCPConn, buf []byte) (bool, error) {
	if dst.ivSent == nil {
		dst.ivSent = new(bool)
	}
	readCount, errRead := conn.Read(buf)
	if errRead != nil {
		if errRead != io.EOF {
			return false, errRead
		}
		return false, nil
	}
	if readCount > 0 {
		writeCount, errWrite := (&TCPConn{
			L:         dst.L,
			R:         dst.R,
			Reader:    bufio.NewReader(dst),
			Writer:    dst,
			Closer:    dst,
			Encryptor: dst.Encryptor,
			ivSent:    dst.ivSent,
		}).ssEncodeWrite(buf[:readCount])
		if errWrite != nil {
			return false, errWrite
		}
		stream.Counter.FlowInIncr(dst.R, uint64(writeCount))
		if readCount != writeCount && writeCount == 0 {
			return false, io.ErrShortWrite
		}
	}
	return true, nil
}
func ssDecodeCopy(conn, dst *TCPConn, buf, iv []byte) (bool, error) {
	readCount, errRead := conn.Read(buf)
	if errRead != nil {
		if errRead != io.EOF {
			return false, errRead
		}
		return false, nil
	}
	if readCount > 0 {
		buf = (*conn.Encryptor).Decode(buf[:readCount], iv)
		writeCount, errWrite := dst.Write(buf)
		if errWrite != nil {
			return false, errWrite
		}
		stream.Counter.FlowOutIncr(conn.R, uint64(writeCount))
		if readCount != writeCount {
			return false, io.ErrShortWrite
		}
	}
	return true, nil
}
