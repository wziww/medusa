package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github/wziww/medusa/log"
	"io"
)

type aesCfb struct {
	password *[]byte
	iv       *[]byte
}

// var _ Encryptor = (*aesCfb)(nil)

// NewAesCfb constructor...
func NewAesCfb(password *[]byte, iv *[]byte) *aesCfb {
	if len(*password) != 16 && len(*password) != 24 && len(*password) != 32 {
		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: password长度必须为16、24或32位"))
		return nil
	}
	if len(*iv) != 16 {
		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: iv长度必须为16位"))
		return nil
	}
	ctr := &aesCfb{password, iv}
	return ctr
}

// Decode ...
func (st *aesCfb) Decode(buf []byte) []byte {

	block, err := aes.NewCipher(*st.password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	if len(buf) < aes.BlockSize {
		log.FMTLog(log.LOGERROR, errors.New("aes_cfb: ciphertext too short"))
		return nil
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(buf, buf)

	return buf
}

// Encode ...
func (st *aesCfb) Encode(buf []byte) []byte {
	block, err := aes.NewCipher(*st.password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(buf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], buf)
	return ciphertext
}
