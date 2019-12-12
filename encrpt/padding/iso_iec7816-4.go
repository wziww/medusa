package encrpt

import "bytes"

// ISO/IEC 7816-4 is identical to the bit padding scheme, applied to a plain text of N bytes.
// This means in practice that the first byte is a mandatory byte valued '80' (Hexadecimal) followed,
// if needed, by 0 to N-1 bytes set to '00', until the end of the block is reached.
// ISO/IEC 7816-4 itself is a communication standard for smart cards containing a file system,
// and in itself does not contain any cryptographic specifications.

// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 80 00 00 00 |
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD DD DD DD 80 |

// ISOIEC7816_4Padding ...
func ISOIEC7816_4Padding(cipherData []byte, blockSize int) []byte {
	// cipherData = append(cipherData, byte(0x80))
	padlen := blockSize - len(cipherData)%blockSize
	padding := []byte{0x80}
	if padlen > 1 {
		padding = append(padding, bytes.Repeat([]byte{0}, padlen-1)...)
	}
	return append(cipherData, padding...)
}

// ISOIEC7816_4UnPadding ...
func ISOIEC7816_4UnPadding(rawData []byte, blockSize int) ([]byte, error) {
	trimZero := bytes.TrimRightFunc(rawData, func(r rune) bool {
		return r == rune(0)
	})
	return trimZero[:len(trimZero)-1], nil
}
