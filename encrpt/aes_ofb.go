package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github/wziww/medusa/log"
	"io"
)

// AesOfb ...
type AesOfb struct {
	Password    *[]byte
	PaddingMode string
}

var _ Encryptor = (*AesOfb)(nil)

// Decode ...
func (st *AesOfb) Decode(cipherBuf []byte) []byte {
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	iv := cipherBuf[:aes.BlockSize]
	buf := cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(buf, buf)
	return buf
}

// Encode ...
func (st *AesOfb) Encode(plainBuf []byte) []byte {
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)
	cipherBuf := make([]byte, aes.BlockSize+len(plainBuf))
	iv := cipherBuf[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGDEBUG, err)
		// return nil
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(cipherBuf[aes.BlockSize:], plainBuf)
	return cipherBuf

}

// Construct ...
func (st *AesOfb) Construct(name string) interface{}{
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
		return nil
	}
	return st
}