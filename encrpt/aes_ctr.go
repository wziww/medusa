package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"strconv"

	"github/wziww/medusa/log"
	"io"
)

// AesCtr ...
type AesCtr struct {
	Password    *[]byte
	PaddingMode string
	cipherBlock cipher.Block
}

var _ Encryptor = (*AesCtr)(nil)

// Decode ...
func (st *AesCtr) Decode(cipherBuf []byte) []byte {
	iv := cipherBuf[:aes.BlockSize]
	stream := cipher.NewCTR(st.cipherBlock, iv)
	buf := cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream.XORKeyStream(buf, buf)
	return buf
}

// Encode ...
func (st *AesCtr) Encode(plainBuf []byte) []byte {
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)
	cipherBuf := make([]byte, aes.BlockSize+len(plainBuf))
	iv := cipherBuf[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	stream := cipher.NewCTR(st.cipherBlock, iv)
	stream.XORKeyStream(cipherBuf[aes.BlockSize:], plainBuf)
	return cipherBuf
}

// Construct ...
func (st *AesCtr) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-ctr":
		targetKeySize = 16
	case "aes-192-ctr":
		targetKeySize = 24
	case "aes-256-ctr":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: key size is"+strconv.Itoa(len(*st.Password))+"should be "+strconv.Itoa(targetKeySize)))
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
