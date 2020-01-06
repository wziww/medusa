package padding

import (
	"bytes"
	"errors"
)

// between 1 and 8 bytes are always added as padding.
// The block is padded with random bytes and the last byte of the block is set to the number of bytes added
// Example: In the following example the block size is 8 bytes, and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 00 00 00 04 |

// ANSIX923Padding ...
func ANSIX923Padding(cipherData []byte, blockSize int) []byte {
	padLen := blockSize - len(cipherData)%blockSize
	padding := bytes.Repeat([]byte{0}, padLen-1)
	return append(cipherData, append(padding, byte(padLen))...)
}

// ANSIX923UnPadding ...
func ANSIX923UnPadding(rawData []byte, blockSize int) ([]byte, error) {
	rawLen := len(rawData)
	if rawLen == 0 {
		return nil, errors.New("ansix923: Raw data is empty")
	}
	if rawLen%blockSize != 0 {
		return nil, errors.New("ansix923: Raw data is not block-aligned")
	}
	end := rawData[rawLen-1]
	padLen := int(end)
	ref := bytes.Repeat([]byte{byte(0)}, padLen-1)
	ref = append(ref, end)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(rawData, ref) {
		return nil, errors.New("ansix923: Invalid padding")
	}
	return rawData[:rawLen-padLen], nil
}
