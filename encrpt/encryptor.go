package encrpt

// Encryptor ...
type Encryptor interface {
	Decode(buf []byte) []byte
	Encode(buf []byte) []byte
	Construct(name string) interface{}
}

// InitEncrypto ...
func InitEncrypto(password *[]byte, method string, paddingMode string) Encryptor {
	switch method {
	case "aes-128-gcm":
		return &AesGcm{password}
	case "aes-128-cfb":
		return &AesCfb{password, paddingMode}
	case "aes-128-ctr":
		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
			return v.(Encryptor)
		}
	case "aes-192-ctr":
		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
			return v.(Encryptor)
		}
	case "aes-256-ctr":
		if v := (&AesCtr{password, paddingMode}).Construct("aes-128-ctr"); v != nil {
			return v.(Encryptor)
		}
	case "aes-128-cbc":
		if v := (&AesCbc{password, paddingMode}).Construct("aes-128-cbc"); v != nil {
			return v.(Encryptor)
		}
	case "aes-192-cbc":
		if v := (&AesCbc{password, paddingMode}).Construct("aes-192-cbc"); v != nil {
			return v.(Encryptor)
		}
	case "aes-256-cbc":
		if v := (&AesCbc{password, paddingMode}).Construct("aes-256-cbc"); v != nil {
			return v.(Encryptor)
		}
	case "aes-128-ofb":
		if v := (&AesOfb{password, paddingMode}).Construct("aes-128-ofb"); v != nil {
			return v.(Encryptor)
		}
	case "aes-192-ofb":
		if v := (&AesOfb{password, paddingMode}).Construct("aes-192-ofb"); v != nil {
			return v.(Encryptor)
		}
	case "aes-256-ofb":
		if v := (&AesOfb{password, paddingMode}).Construct("aes-256-ofb"); v != nil {
			return v.(Encryptor)
		}
	}
	return nil
}

type encryptInfo struct {
	// name string
	keySize   int
	encryptor interface{}
}

var encryptorMap = map[string]*encryptInfo{
	// "aes-128-cbc": {16, AesCbc},
}

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

// CheckAndGen ...
func CheckAndGen(name string, encryptor *Encryptor) interface{} {
	// if len(*password) != size {
	// 	// log.FMTLog(log.LOGERROR, errors.New("aes_cbc: key size should be "+strconv.Itoa(size)))
	// 	return nil
	// }
	if v := (*encryptor).Construct("aes-256-ofb"); v != nil {
		return v.(Encryptor)
	}
	// if v := (&AesOfb{password, paddingMode}).Construct("aes-256-ofb"); v != nil {
	// 		return v.(Encryptor)
	// 	}
	return nil
}
