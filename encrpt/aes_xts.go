package encrpt

import (
	"crypto/aes"
	"github/wziww/medusa/log"

	"golang.org/x/crypto/xts"
)

// AesXts ...
type AesXts struct {
	Password *[]byte
	Sector   *uint64
}

var _ Encryptor = (*AesXts)(nil)

// Decode ...
func (st *AesXts) Decode(cipherBuf []byte) []byte {
	cipher, err := xts.NewCipher(aes.NewCipher, *st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	buf := make([]byte, len(cipherBuf))
	cipher.Encrypt(buf, cipherBuf, *st.Sector)
	return buf
}

// Encode ...
func (st *AesXts) Encode(plainBuf []byte) []byte {

	cipher, err := xts.NewCipher(aes.NewCipher, *st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	cipherBuf := make([]byte, len(plainBuf))
	cipher.Encrypt(cipherBuf, plainBuf, *st.Sector)
	return cipherBuf
}

// Construct ...
func (st *AesXts) Construct(name string) interface{}{
		var targetKeySize int
	switch name {
	case "aes-128-xts":
		targetKeySize = 16
	case "aes-192-xts":
		targetKeySize = 24
	case "aes-256-xts":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		return nil
	}
	return st
}