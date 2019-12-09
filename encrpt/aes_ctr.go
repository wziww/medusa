package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github/wziww/medusa/log"
)

// Aes128ctr ...
// type Aes128ctr struct {
// 	Password *[]byte
// 	IV *[]byte
// }

type aesCtr struct {
	password *[]byte
	iv       *[]byte
}

// var _ Encryptor = (*aesCtr)(nil)

// NewAesCtr constructor...
func NewAesCtr(password *[]byte, iv *[]byte) *aesCtr {
	if len(*password) != 16 && len(*password) != 24 && len(*password) != 32 {
		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: password长度必须为16、24或32位"))
		return nil
	}
	if len(*iv) != 16 {
		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: iv长度必须为16位"))
		return nil
	}
	ctr := &aesCtr{password, iv}
	return ctr
}

func (st *aesCtr) Decode(buf []byte) []byte {
	block, err := aes.NewCipher(*st.password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	stream := cipher.NewCTR(block, *st.iv)
	plainBuf := make([]byte, len(buf))
	stream.XORKeyStream(plainBuf, buf)
	return plainBuf
}

func (st *aesCtr) Encode(buf []byte) []byte {
	block, err := aes.NewCipher(*st.password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	stream := cipher.NewCTR(block, *st.iv)
	cipherBuf := make([]byte, len(buf))
	stream.XORKeyStream(cipherBuf, buf)
	return cipherBuf
}
