package encrpt

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte) []byte
	Encode(buf []byte) []byte
	Construct(name string) interface{}
}

type encryptInfo struct {
	keySize   int
	encryptor Encryptor
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
	var encryptorMap = map[string]*encryptInfo{
		"aes-128-cbc": {16, &AesCbc{password, paddingMode, nil}},
		"aes-192-cbc": {24, &AesCbc{password, paddingMode, nil}},
		"aes-256-cbc": {32, &AesCbc{password, paddingMode, nil}},
		"aes-128-cfb": {16, &AesCfb{password, paddingMode, nil}},
		"aes-192-cfb": {24, &AesCfb{password, paddingMode, nil}},
		"aes-256-cfb": {32, &AesCfb{password, paddingMode, nil}},
		"aes-128-ctr": {16, &AesCtr{password, paddingMode, nil}},
		"aes-192-ctr": {24, &AesCtr{password, paddingMode, nil}},
		"aes-256-ctr": {32, &AesCtr{password, paddingMode, nil}},
		"aes-128-gcm": {16, &AesGcm{password, paddingMode, nil}},
		"aes-192-gcm": {24, &AesGcm{password, paddingMode, nil}},
		"aes-256-gcm": {32, &AesGcm{password, paddingMode, nil}},
		"aes-128-ofb": {16, &AesOfb{password, paddingMode, nil}},
		"aes-192-ofb": {24, &AesOfb{password, paddingMode, nil}},
		"aes-256-ofb": {32, &AesOfb{password, paddingMode, nil}},
	}
	enc := encryptorMap[method]
	return CheckAndGen(method, &(enc.encryptor))
}

// CheckAndGen ...
func CheckAndGen(name string, encryptor *Encryptor) Encryptor {
	if v := (*encryptor).Construct(name); v != nil {
		return v.(Encryptor)
	}
	return nil
}
