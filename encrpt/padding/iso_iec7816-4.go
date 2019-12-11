package encrpt

// specifies that the padding should be done at the end of that last block with random bytes, 
// and the padding boundary should be specified by the last byte.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 81 A6 23 04 |
//


// ISOIEC7816_4Padding ...
func ISOIEC7816_4Padding(ciphertext []byte, blockSize int) []byte {
	return nil
}

// ISOIEC7816_4UnPadding ...
func ISOIEC7816_4UnPadding(origData []byte) []byte {
	return nil
}
