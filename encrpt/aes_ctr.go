package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	// "errors"
	"github/wziww/medusa/log"
	"io"
)

// AesCtr ...
type AesCtr struct {
	Password    *[]byte
	PaddingMode string
}

var _ Encryptor = (*AesCtr)(nil)

// NewAesCtr constructor...
// func NewAesCtr(password *[]byte, iv *[]byte) *aesCtr {
// 	if len(*password) != 16 && len(*password) != 24 && len(*password) != 32 {
// 		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: password长度必须为16、24或32位"))
// 		return nil
// 	}
// 	if len(*iv) != 16 {
// 		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: iv长度必须为16位"))
// 		return nil
// 	}
// 	ctr := &aesCtr{password, iv}
// 	return ctr
// }

// Decode ...
func (st *AesCtr) Decode(cipherBuf []byte) []byte {
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	iv := cipherBuf[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	// plainBuf := make([]byte, len(cipherBuf))
	buf := cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream.XORKeyStream(buf, buf)
	return buf
}

// Encode ...
func (st *AesCtr) Encode(plainBuf []byte) []byte {
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
	stream := cipher.NewCTR(block, iv)
	// cipherBuf := make([]byte, len(plainBuf))
	stream.XORKeyStream(cipherBuf[aes.BlockSize:], plainBuf)
	return cipherBuf
}
