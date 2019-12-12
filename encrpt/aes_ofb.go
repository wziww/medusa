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
	Password *[]byte
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
