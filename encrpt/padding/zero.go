package encrpt

// All the bytes that are required to be padded are padded with zero.
// Example: In the following example the block size is 8 bytes and padding is required for 4 bytes
// ... | DD DD DD DD DD DD DD DD | DD DD DD DD 00 00 00 00 |

// ZeroPadding ...
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	return nil
}

// ZeroUnPadding ...
func ZeroUnPadding(origData []byte) []byte {
	return nil
}
