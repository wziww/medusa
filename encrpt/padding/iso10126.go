package encrpt

import (
	"crypto/rand"
	"errors"
	"github/wziww/medusa/log"
	"io"
)

// ISO10126Padding ...
func ISO10126Padding(cipherData []byte, blockSize int) []byte {
	padLen := blockSize - len(cipherData)%blockSize
	padding := make([]byte, padLen-1)
	if _, err := io.ReadFull(rand.Reader, padding); err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	return append(cipherData, append(padding, byte(padLen))...)
}

// ISO10126UnPadding ...
func ISO10126UnPadding(rawData []byte, blockSize int) ([]byte, error) {
	rawLen := len(rawData)
	if rawLen == 0 {
		return nil, errors.New("iso10126: Raw data is empty")
	}
	if rawLen%blockSize != 0 {
		return nil, errors.New("iso10126: Raw data is not block-aligned")
	}
	padLen := int(rawData[rawLen-1])
	// random byte, ignore check
	if padLen > blockSize || padLen == 0 {
		return nil, errors.New("iso10126: Invalid padding")
	}
	return rawData[:rawLen-padLen], nil
}
