package medusa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github/wziww/medusa/log"
	"io"
)

// Aes128gcm ...
type Aes128gcm struct {
	Password *[]byte
}

var _ Encryptor = (*Aes128gcm)(nil)

// Decode ...
func (st *Aes128gcm) Decode(buf []byte) []byte {
	nonce, _ := hex.DecodeString("000000000000000000000000") //加密用的nonce
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	plaintext, err := aesgcm.Open(nil, nonce, buf, nil)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return plaintext
}

// Encode ...
func (st *Aes128gcm) Encode(buf []byte) []byte {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	nonce := make([]byte, 12)
	if false {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			log.FMTLog(log.LOGERROR, err)
		}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	return aesgcm.Seal(nil, nonce, buf, nil)
}
