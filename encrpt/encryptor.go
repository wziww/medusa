package encrpt

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte) []byte
	Encode(buf []byte) []byte
	Construct(name string) interface{}
}

// InitEncrypto ...
// func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
// 	switch method {
// 	case "aes-128-gcm":
// 		return &AesGcm{password,paddingMode}
// 	case "aes-128-cfb":
// 		return &AesCfb{password, paddingMode}
// 	case "aes-128-ctr":
// 		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-192-ctr":
// 		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-256-ctr":
// 		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-128-cbc":
// 		if v := (&AesCbc{password, paddingMode}).Construct("aes-128-cbc"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-192-cbc":
// 		if v := (&AesCbc{password, paddingMode}).Construct("aes-192-cbc"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-256-cbc":
// 		if v := (&AesCbc{password, paddingMode}).Construct("aes-256-cbc"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-128-ofb":
// 		if v := (&AesOfb{password, paddingMode}).Construct("aes-128-ofb"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-192-ofb":
// 		if v := (&AesOfb{password, paddingMode}).Construct("aes-192-ofb"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	case "aes-256-ofb":
// 		if v := (&AesOfb{password, paddingMode}).Construct("aes-256-ofb"); v != nil {
// 			return v.(Encryptor)
// 		}
// 	}
// 	return nil
// }

// InitEncrypto ...
// func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
// 	switch method {
// 	case "aes-128-gcm":
// 		return &Aes128gcm{password}
// 	case "aes-128-cfb":
// 		return &AesCfb{password, paddingMode}
// 	case "aes-128-ctr":
// 		return &AesCtr{password, paddingMode}
// 	case "aes-128-cbc":
// 		return &AesCbc{password, paddingMode}
// 	case "aes-128-ofb":
// 		return &AesOfb{password, paddingMode}
// 	}
// 	return nil
// }

type encryptInfo struct {
	keySize   int
	encryptor Encryptor
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
	var encryptorMap = map[string]*encryptInfo{
		"aes-128-cbc": {16, &AesCbc{password, paddingMode,nil}},
		"aes-192-cbc": {24, &AesCbc{password, paddingMode,nil}},
		"aes-256-cbc": {32, &AesCbc{password, paddingMode,nil}},
		"aes-128-cfb": {16, &AesCfb{password, paddingMode}},
		"aes-192-cfb": {24, &AesCfb{password, paddingMode}},
		"aes-256-cfb": {32, &AesCfb{password, paddingMode}},
		"aes-128-ctr": {16, &AesCtr{password, paddingMode}},
		"aes-192-ctr": {24, &AesCtr{password, paddingMode}},
		"aes-256-ctr": {32, &AesCtr{password, paddingMode}},
		"aes-128-gcm": {16, &AesGcm{password, paddingMode}},
		"aes-192-gcm": {24, &AesGcm{password, paddingMode}},
		"aes-256-gcm": {32, &AesGcm{password, paddingMode}},
		"aes-128-ofb": {16, &AesOfb{password, paddingMode}},
		"aes-192-ofb": {24, &AesOfb{password, paddingMode}},
		"aes-256-ofb": {32, &AesOfb{password, paddingMode}},
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
