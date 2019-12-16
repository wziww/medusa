package encrpt

import (
	"errors"
	"github/wziww/medusa/log"
)

// EncrptFactory ...
type EncrptFactory interface {
	GenEncrptor(password *[]byte, paddingMode string, t string) Encrptor
}

// Encrptor ...
type Encrptor interface {
	// Decode(buf []byte) []byte
	// Encode(buf []byte) []byte
	Decode(buf string) string
	Encode(buf string) string
}

// InitEncrptor ...
// func InitEncrptor(password *[]byte, method string, paddingMode string) Encryptor {
// 	switch method {
// 	case "aes-128-gcm":
// 		return &Aes128gcm{password}
// 	case "aes-128-cfb":
// 		return &AesCfb{password, paddingMode}
// 	case "aes-128-ctr":
// 		return &AesCtr{password, paddingMode}
// 	case "aes-128-cbc":
// 		return &AesCbc{password, paddingMode}
// 	case "aes-128-ofb":
// 		return &AesOfb{password, paddingMode}
// 	}
// 	return nil
// }

// AesEncrptFactory ...
type AesEncrptFactory struct{}

// AesCbcEncrptFactory ...
type AesCbcEncrptFactory struct{}

// GenEncrptor ...
func (t *AesCbcEncrptFactory) GenEncrptor(password *[]byte, paddingMode string, s string) Encrptor {
	switch s {
	case "aes-128-cbc":
		{
			if len(*password) != 16 {
				log.FMTLog(log.LOGERROR, errors.New("aes_cbc: len(key) should be 16"))
				return nil
			}
			return &Aes128CbcEncrptor{&AesCbcEncrptor{password, paddingMode}}
		}
	case "aes-192-cbc":
		{
			if len(*password) != 24 {
				log.FMTLog(log.LOGERROR, errors.New("aes_cbc: len(key) should be 24"))
				return nil
			}
			return &Aes128CbcEncrptor{&AesCbcEncrptor{password, paddingMode}}
		}
	case "aes-256-cbc":
		{
			if len(*password) != 32 {
				log.FMTLog(log.LOGERROR, errors.New("aes_cbc: len(key) should be 32"))
				return nil
			}
			return &Aes128CbcEncrptor{&AesCbcEncrptor{password, paddingMode}}
		}
	default:
		return nil
	}
}

// AesCbcEncrptor ...
type AesCbcEncrptor struct {
	Password    *[]byte
	PaddingMode string
}

// Decode ...
func (st *AesCbcEncrptor) Decode(s string) string {
	return "lalala"
}

// Encode ...
func (st *AesCbcEncrptor) Encode(s string) string {
	return "papapa"
}

// Aes128CbcEncrptor ...
type Aes128CbcEncrptor struct {
	*AesCbcEncrptor
}

// Decode ...
func (st *Aes128CbcEncrptor) Decode(s string) string {
	return (*st.AesCbcEncrptor).Decode(s)
}

// Encode ...
func (st *Aes128CbcEncrptor) Encode(s string) string {
	return (*st.AesCbcEncrptor).Encode(s)
}

// Aes192CbcEncrptor ...
type Aes192CbcEncrptor struct {
	*AesCbcEncrptor
}

// Aes256CbcEncrptor ...
type Aes256CbcEncrptor struct {
	*AesCbcEncrptor
}

// Decode ...
func (st *Aes192CbcEncrptor) Decode(s string) string {
	return (*st.AesCbcEncrptor).Decode(s)
}

// Encode ...
func (st *Aes192CbcEncrptor) Encode(s string) string {
	return (*st.AesCbcEncrptor).Encode(s)
}

// Decode ...
func (st *Aes256CbcEncrptor) Decode(s string) string {
	return (*st.AesCbcEncrptor).Decode(s)
}

// Encode ...
func (st *Aes256CbcEncrptor) Encode(s string) string {
	return (*st.AesCbcEncrptor).Encode(s)
}
