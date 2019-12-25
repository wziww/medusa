package medusa

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github/wziww/medusa/encrpt"
	"github/wziww/medusa/log"
	"github/wziww/medusa/stream"
	"io"
	"sync"
	"sync/atomic"
)

const tag = 0x80

var bufSize int64 = 1 << 10 //1kb

var bp sync.Pool

const maxConsecutiveEmptyReads = 100

func init() {
	bp.New = func() interface{} {
		b := make([]byte, 8)
		return &b
	}
}

func btsPoolGet() *[]byte {
	return bp.Get().(*[]byte)
}

func btsPoolPut(b *[]byte) {
	bp.Put(b)
}

// TCPConn ...
type TCPConn struct {
	L string // LocalAddr
	R string
	*bufio.Reader
	io.Closer
	io.Writer
	Encryptor encrpt.Encryptor
	MutexR    sync.Mutex
	MutexW    sync.Mutex
}

// DecodeRead ...
func (conn *TCPConn) DecodeRead() (n int, buf []byte, err error) {
	conn.MutexR.Lock()
	defer conn.MutexR.Unlock()
	// /**
	//   +----+-----+-------+------+----------+----------+
	//   |LEN | 								DATA 										 |
	//   +----+-----+-------+------+----------+----------+
	//   | 8  | 								 x  									   |
	//   +----+-----+-------+------+----------+----------+
	// */
	var l int64
	// binary.Read(conn, binary.BigEndian, &l)
	b := btsPoolGet()
	defer btsPoolPut(b)
	var readN int
	readN, err = conn.Read(*b)
	if err != nil {
		log.FMTLog(log.LOGDEBUG, err)
		return
	}
	if readN < 8 {
		// 非法协议
		return
	}
	l = int64(binary.BigEndian.Uint64((*b)[:8]))
	if l <= 0 || l > (bufSize<<10) {
		return
	}
	atomic.AddUint64(stream.FlowIn, uint64(l))
	// /**
	//   +----+-----+-------+------+----------+----------+
	//   |LEN | 								DATA 										 |
	//   +----+-----+-------+------+----------+----------+
	//   | 8  | 					   readN - 8								   |
	//   +----+-----+-------+------+----------+----------+
	//   +----+-----+-------+------+----------+----------+
	//   |         							DATA 										 |
	//   +----+-----+-------+------+----------+----------+
	//   |    							l - readN + 8 		 				   |
	//   +----+-----+-------+------+----------+----------+
	//   +----+-----+-------+------+----------+----------+
	//   |         						DATAALL										 |
	//   +----+-----+-------+------+----------+----------+
	//   |    							l             		 				   |
	//   +----+-----+-------+------+----------+----------+
	//
	// */
	data := make([]byte, l)
	var rm int
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		rn, _ := conn.Read(data[rm:])
		if rn < 0 {
			return 0, nil, errors.New("bufio: reader returned negative count from Read")
		}
		rm += rn
		if int64(rm) == l {
			res := conn.Encryptor.Decode(data)
			if res != nil {
				buf = res
				n = len(res)
				return
			}
			return 0, nil, errors.New("DecodeRead error")
		}
	}
	return 0, nil, io.ErrNoProgress
}

//EncodeWrite ...
func (conn *TCPConn) EncodeWrite(buf []byte) (n int, err error) {
	buf = conn.Encryptor.Encode(buf)
	if buf != nil {
		// /**
		//   +----+-----+-------+------+----------+----------+
		//   |LEN | 								DATA 										 |
		//   +----+-----+-------+------+----------+----------+
		//   | 8  | 								 x  									   |
		//   +----+-----+-------+------+----------+----------+
		// */
		var l int64 = int64(len(buf))
		atomic.AddUint64(stream.FlowOut, uint64(l))
		conn.MutexW.Lock()
		defer conn.MutexW.Unlock()
		binary.Write(conn, binary.BigEndian, l)
		return conn.Write(buf)
	}
	return
}

// DecodeCopy 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) DecodeCopy(dst *TCPConn) error {
	for {
		c, err := decodeCopy(conn, dst)
		if !c {
			return err
		}
	}
}

func decodeCopy(conn, dst *TCPConn) (bool, error) {
	readCount, buf, errRead := conn.DecodeRead()
	if errRead != nil {
		if errRead != io.EOF {
			return false, errRead
		}
		return false, nil
	}
	if readCount > 0 {
		writeCount, errWrite := dst.Write(buf[:readCount])
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

// EncodeCopy 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) EncodeCopy(dst *TCPConn) error {
	buf := make([]byte, bufSize)
	for {
		c, err := encodeCopy(conn, dst, buf)
		if !c {
			return err
		}
	}
}
func encodeCopy(conn, dst *TCPConn, buf []byte) (bool, error) {
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
			Encryptor: conn.Encryptor,
		}).EncodeWrite(buf[:readCount])
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
