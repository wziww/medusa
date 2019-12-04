package medusa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github/wziww/medusa/config"
	"github/wziww/medusa/log"
	"io"
)

const (
	bufSize = 1024
)

// TCPConn ...
type TCPConn struct {
	io.ReadWriteCloser
}

func decode(buf []byte) []byte {
	key := []byte(config.C.Base.Password)
	nonce, _ := hex.DecodeString("000000000000000000000000") //加密用的nonce
	block, err := aes.NewCipher(key)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	plaintext, err := aesgcm.Open(nil, nonce, buf, nil)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return plaintext
}

// DecodeRead ...
func (conn *TCPConn) DecodeRead(buf []byte) (n int, err error) {
	n, err = conn.Read(buf)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return
	}
	if res := decode(buf[:n]); res != nil {
		buf = res
		return
	}
	return
}
func encode(buf []byte) []byte {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := []byte(config.C.Base.Password)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	nonce := make([]byte, 12)
	if false {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			log.FMTLog(log.LOGERROR, err)
		}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	return aesgcm.Seal(nil, nonce, buf, nil)
}

//EncodeWrite ...
func (conn *TCPConn) EncodeWrite(buf []byte) (n int, err error) {
	buf = encode(buf)
	if buf != nil {
		return conn.Write(buf)
	}
	return
}

// DecodeCopy 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (conn *TCPConn) DecodeCopy(dst io.Writer) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := conn.DecodeRead(buf)
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
			}).EncodeWrite(buf[:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
