package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github/wziww/medusa/log"
	"io"
)

// AesCbc ...
type AesCbc struct {
	Password    *[]byte
	PaddingMode string
}

var _ Encryptor = (*AesCbc)(nil)

// Decode ...
func (st *AesCbc) Decode(cipherBuf []byte) []byte {
	// block, err := aes.NewCipher(*st.Password)
	block, err := GetSingleCipher(st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	if len(cipherBuf) < aes.BlockSize {
		log.FMTLog(log.LOGERROR, errors.New("aes_cbc: cipherBuf too short"))
		return nil
	}
	iv := cipherBuf[:aes.BlockSize]
	cipherBuf = cipherBuf[aes.BlockSize:]
	if len(cipherBuf)%aes.BlockSize != 0 {
		log.FMTLog(log.LOGERROR, errors.New("aes_cbc: cipherBuf is not a multiple of the block size"))
		return nil
	}

	blockMode := cipher.NewCBCDecrypter(*block, iv)
	blockMode.CryptBlocks(cipherBuf, cipherBuf)
	// unpad
	cipherBuf, _ = HandleUnPadding(st.PaddingMode)(cipherBuf, aes.BlockSize)
	return cipherBuf
}

// Encode ...
func (st *AesCbc) Encode(plainBuf []byte) []byte {
	if len(plainBuf)%aes.BlockSize != 0 {
		log.FMTLog(log.LOGERROR, errors.New("aes_cbc: plainBuf is not a multiple of the block size"))
		return nil
	}
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)
	// block, err := aes.NewCipher(*st.Password)
	block, err := GetSingleCipher(st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	cipherBuf := make([]byte, aes.BlockSize+len(plainBuf))
	iv := cipherBuf[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGERROR, err)
	}
	blockMode := cipher.NewCBCEncrypter(*block, iv)
	blockMode.CryptBlocks(cipherBuf[aes.BlockSize:], plainBuf)
	return cipherBuf
}

// Construct ...
func (st *AesCbc) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-cbc":
		targetKeySize = 16
	case "aes-192-cbc":
		targetKeySize = 24
	case "aes-256-cbc":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		// log.FMTLog(log.LOGERROR, errors.New("aes_cbc: key size should be "+strconv.Itoa(size)))
		return nil
	}
	return st
}
