package encrypt

import (
	"github/wziww/medusa/config"
	"testing"
)

func TestMain(m *testing.M) {
	config.Init()
	m.Run()
}

func TestGetpassword(t *testing.T) {
	a, b := getpassword([]byte("abcdascsdjfgw"), 16, 32)
	if string(a) !=
		string([]byte{219, 123, 80, 239, 101, 230, 14, 137, 225, 203, 178, 179, 234, 83, 101, 241}) || string(b) !=
		string([]byte{46, 132, 232, 48, 185, 224, 186, 29, 41, 210, 80, 221, 180, 199, 235, 56, 156, 76, 10, 110, 60, 224, 213, 212, 76, 4, 128, 194, 194, 14, 236, 35}) {
		t.Fatal("getpassword error")
	}
}
func TestGetRandString(t *testing.T) {
	if 33 != len(GetRandString(33)) {
		t.Fatal("GetRandString error")
	}
}
func TestInitEncrypto(t *testing.T) {
	password := []byte("xscuagi2qeg")
	z := InitEncrypto(&password, "", "", nil)
	if z != nil {
		t.Fatal("InitEncrypto error")
	}
	z = InitEncrypto(&password, "aes-256-cfb", "", nil)
	if z == nil {
		t.Fatal("InitEncrypto error")
	}
}
