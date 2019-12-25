package  padding

import (
	"bytes"
	"errors"
)

// The value of each added byte is the number of bytes that are added, i.e. N bytes, each of value N are added.
// The number of bytes added will depend on the block boundary to which the message needs to be extended.
// This padding method (as well as the previous two) is well-defined if and only if N is less than 256.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 04 04 04 04 |

// PKCS7Padding ...
func PKCS7Padding(cipherData []byte, blockSize int) []byte {
	// if blockSize < 0 || blockSize > 256 {
	// 	return nil, errors.New("pkcs7: Invalid block size " + strconv.Itoa(blockSize))
	// }
	padLen := blockSize - len(cipherData)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(cipherData, padding...)

}

// PKCS7UnPadding ...
func PKCS7UnPadding(rawData []byte, blockSize int) ([]byte, error) {
	rawLen := len(rawData)
	if rawLen == 0 {
		return nil, errors.New("pkcs7: Raw data is empty")
	}
	if rawLen%blockSize != 0 {
		return nil, errors.New("pkcs7: Raw data is not block-aligned")
	}
	padLen := int(rawData[rawLen-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(rawData, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return rawData[:rawLen-padLen], nil

}
