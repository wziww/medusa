package encrypt

import (
	"bytes"
	"crypto/md5"
	"math/rand"
	"time"
)

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte, iv ...[]byte) []byte
	Encode(buf []byte) []byte
	Construct(name string, iv []byte) interface{}
	Ivlen() int
	GetIv() []byte
}

func getpassword(password []byte, keyLen, ivLen int) ([]byte, []byte) {
	var m [][]byte
	i := 0
	for len(m) < (keyLen + ivLen) {
		data := password
		if i > 0 {
			data = []byte(string(m[i-1]) + string(password))
		}
		digst := md5.Sum(data)
		m = append(m, digst[:])
		i++
	}
	data := bytes.Join(m, []byte{})
	return data[:keyLen], data[keyLen : keyLen+ivLen]
}

// GetRandString returns randominzed string with given length
func GetRandString(length int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string, paddingMode string, iv []byte) Encryptor {
	z := (&CipherS{
		Password: password,
		Method:   method,
	}).Construct(method, iv)
	if z == nil {
		return nil
	}
	return z.(Encryptor)
}
