package medusa

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte) []byte
	Encode(buf []byte) []byte
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string) Encryptor {
	switch method {
	case "aes-128-gcm":
		return &Aes128gcm{
			Password: password,
		}
	}
	return nil
}
