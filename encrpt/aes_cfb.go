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

// AesCfb ...
type AesCfb struct {
	Password    *[]byte
	PaddingMode string
	cipherBlock cipher.Block
}

var _ Encryptor = (*AesCfb)(nil)

// Decode ...
func (st *AesCfb) Decode(cipherBuf []byte) []byte {
	if len(cipherBuf) < aes.BlockSize {
		log.FMTLog(log.LOGERROR, errors.New("aes_cfb: ciphertext too short"))
		return nil
	}
	iv := cipherBuf[:aes.BlockSize]
	var buf = cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream := cipher.NewCFBDecrypter(st.cipherBlock, iv)
	stream.XORKeyStream(buf, buf)

	return buf
}

// Encode ...
func (st *AesCfb) Encode(plainBuf []byte) []byte {
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plainBuf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	stream := cipher.NewCFBEncrypter(st.cipherBlock, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainBuf)
	return ciphertext
}

// Construct ...
func (st *AesCfb) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-cfb":
		targetKeySize = 16
	case "aes-192-cfb":
		targetKeySize = 24
	case "aes-256-cfb":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		log.FMTLog(log.LOGERROR, errors.New("aes_cfb: key size is"+strconv.Itoa(len(*st.Password))+"should be "+strconv.Itoa(targetKeySize)))
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
