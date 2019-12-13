package padding

import "bytes"

// All the bytes that are required to be padded are padded with zero.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 00 00 00 00 |
// WARN: 结尾为0的原始数据进行零填充，可能会不可逆

// ZeroPadding ...
func ZeroPadding(cipherData []byte, blockSize int) []byte {
	padlen := blockSize - len(cipherData)%blockSize
	padding := bytes.Repeat([]byte{0}, padlen)
	return append(cipherData, padding...)
}

// ZeroUnPadding ...
func ZeroUnPadding(rawData []byte) ([]byte, error) {
	return bytes.TrimRightFunc(rawData, func(r rune) bool {
		return r == rune(0)
	}), nil
}
