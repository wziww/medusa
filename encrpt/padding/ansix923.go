package encrpt

// between 1 and 8 bytes are always added as padding. 
// The block is padded with random bytes and the last byte of the block is set to the number of bytes added
// Example: In the following example the block size is 8 bytes, and padding is required for 4 bytes 
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 00 00 00 04 |

// ANSIX923Padding ...
func ANSIX923Padding(ciphertext []byte, blockSize int) []byte{
	return nil
}

// ANSIX923UnPadding ...
func ANSIX923UnPadding(origData []byte) []byte{
	return nil
}