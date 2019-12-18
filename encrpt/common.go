package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"sync"
)

var once sync.Once
var cipherBlock *cipher.Block

// GetSingleCipher ...
func GetSingleCipher(key *[]byte) (*cipher.Block, error) {
	var errors error = nil
	once.Do(func() {
		block, err := aes.NewCipher(*key)
		if err != nil {
			errors = err
		}
		cipherBlock = &block
	})
	return cipherBlock, errors
}
