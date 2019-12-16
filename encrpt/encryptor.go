package encrpt

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte) []byte
	Encode(buf []byte) []byte
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
	switch method {
	case "aes-128-gcm":
		return &Aes128gcm{password}
	case "aes-128-cfb":
		return &AesCfb{password, paddingMode}
	case "aes-128-ctr":
		return &AesCtr{password, paddingMode}
	case "aes-128-cbc":
		return &AesCbc{password, paddingMode}
	case "aes-128-ofb":
		return &AesOfb{password, paddingMode}
	}
	return nil
}
