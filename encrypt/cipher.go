package encrypt

import (
	"crypto/cipher"
	"github/wziww/medusa/log"

	"github.com/spacemonkeygo/openssl"
)

// CipherS ...
type CipherS struct {
	Password            *[]byte
	Method              string
	PaddingMode         string
	cipherBlock         cipher.Block
	encodeiv            []byte
	encryptionCipherCtx *openssl.EncryptionCipherCtx
	decryptionCipherCtx *openssl.DecryptionCipherCtx
}

var _ Encryptor = (*CipherS)(nil)

// GetIv ...
func (st *CipherS) GetIv() []byte {
	return st.encodeiv
}

// Ivlen ...
func (st *CipherS) Ivlen() int {
	return len(st.encodeiv)
}

// Decode ...
func (st *CipherS) Decode(cipherBuf []byte, iv ...[]byte) []byte {
	if st.decryptionCipherCtx == nil {
		dcipher, dcipherError := openssl.GetCipherByName(st.Method)
		if dcipherError != nil {
			log.FMTLog(log.LOGERROR, dcipherError)
		}
		var ivv []byte
		if len(iv) > 0 && len(iv[0]) == st.Ivlen() {
			ivv = iv[0]
		} else {
			ivv = st.GetIv()
		}
		dCtx, dCtxError := openssl.NewDecryptionCipherCtx(dcipher, nil, *st.Password, ivv)
		if dCtxError != nil {
			log.FMTLog(log.LOGERROR, dCtxError)
		}
		st.decryptionCipherCtx = &dCtx
	}
	cipherbytes, err := (*st.decryptionCipherCtx).DecryptUpdate(cipherBuf)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return cipherbytes
}

// Encode ...
func (st *CipherS) Encode(plainBuf []byte) []byte {
	cipherbytes, err := (*st.encryptionCipherCtx).EncryptUpdate(plainBuf)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return cipherbytes
}

// Construct ...
func (st *CipherS) Construct(name string, iv []byte) interface{} {
	cipher, cipherError := openssl.GetCipherByName(st.Method)
	if cipherError != nil {
		log.FMTLog(log.LOGERROR, cipherError)
		return nil
	}
	st.Method = name
	password, ivv := getpassword(*st.Password, cipher.KeySize(), cipher.IVSize())
	st.Password = &password
	if iv == nil || len(iv) == 0 {
		iv = ivv
	}
	st.encodeiv = iv
	eCtx, eCtxError := openssl.NewEncryptionCipherCtx(cipher, nil, password, iv)
	if eCtxError != nil {
		log.FMTLog(log.LOGERROR, eCtxError)
		return nil
	}
	st.encryptionCipherCtx = &eCtx
	return st
}
