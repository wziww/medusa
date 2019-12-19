package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github/wziww/medusa/log"
	"io"
	"strconv"
)

// AesOfb ...
type AesOfb struct {
	Password    *[]byte
	PaddingMode string
	cipherBlock cipher.Block
}

var _ Encryptor = (*AesOfb)(nil)

// Decode ...
func (st *AesOfb) Decode(cipherBuf []byte) []byte {
	iv := cipherBuf[:aes.BlockSize]
	buf := cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream := cipher.NewOFB(st.cipherBlock, iv)
	stream.XORKeyStream(buf, buf)
	return buf
}

// Encode ...
func (st *AesOfb) Encode(plainBuf []byte) []byte {
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)
	cipherBuf := make([]byte, aes.BlockSize+len(plainBuf))
	iv := cipherBuf[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGDEBUG, err)
		// return nil
	}

	stream := cipher.NewOFB(st.cipherBlock, iv)
	stream.XORKeyStream(cipherBuf[aes.BlockSize:], plainBuf)
	return cipherBuf

}

// Construct ...
func (st *AesOfb) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-ofb":
		targetKeySize = 16
	case "aes-192-ofb":
		targetKeySize = 24
	case "aes-256-ofb":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		log.FMTLog(log.LOGERROR, errors.New("aes_ofb: key size is"+strconv.Itoa(len(*st.Password))+"should be "+strconv.Itoa(targetKeySize)))
		return nil
	}
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	st.cipherBlock = block
	return st
}
