package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github/wziww/medusa/log"
	"strconv"
)

// AesGcm ...
type AesGcm struct {
	Password    *[]byte
	PaddingMode string
	cipherBlock cipher.Block
}

var _ Encryptor = (*AesGcm)(nil)

// Decode ...
func (st *AesGcm) Decode(buf []byte) []byte {
	nonce, _ := hex.DecodeString("000000000000000000000000") //加密用的nonce
	aesgcm, err := cipher.NewGCM(st.cipherBlock)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	plaintext, err := aesgcm.Open(nil, nonce, buf, nil)
	if err != nil {
		log.FMTLog(log.LOGDEBUG, err)
		return nil
	}
	// unpad
	plaintext, _ = HandleUnPadding(st.PaddingMode)(plaintext, aes.BlockSize)
	return plaintext
}

// Encode ...
func (st *AesGcm) Encode(buf []byte) []byte {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	nonce := make([]byte, 12)

	aesgcm, err := cipher.NewGCM(st.cipherBlock)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	// pad
	buf = HandlePadding(st.PaddingMode)(buf, aes.BlockSize)
	return aesgcm.Seal(nil, nonce, buf, nil)
}

// Construct ...
func (st *AesGcm) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-gcm":
		targetKeySize = 16
	case "aes-192-gcm":
		targetKeySize = 24
	case "aes-256-gcm":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		log.FMTLog(log.LOGERROR, errors.New("aes_gcm: key size is"+strconv.Itoa(len(*st.Password))+"should be "+strconv.Itoa(targetKeySize)))
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
