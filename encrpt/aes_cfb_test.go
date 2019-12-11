package encrpt

import (
	"fmt"
	"testing"
)

var passwordCfb []byte = []byte("AES256Key-32Characters1234567890")
var aesobjCfb *AesCfb = &AesCfb{
	Password: &passwordCfb,
}

func TestStringCfb(t *testing.T) {
	s := "hello world!"
	sd := aesobjCfb.Encode([]byte(s))
	fmt.Printf("sd---%s\n", sd)
	s2 := aesobjCfb.Decode(sd)
	fmt.Printf("s2---%s\n", s2)
	if s != string(s2) {
		t.Fatal(s, "!=", string(s2), "fail to encode and decode")
	}
}

func TestBytesCfb(t *testing.T) {
	s := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	sd := aesobjCfb.Encode([]byte(s))
	s2 := aesobjCfb.Decode(sd)
	for i := range s {
		if s[i] != s2[i] {
			t.Fatal(s, "!=", string(s2), "fail to encode and decode")
			return
		}
	}
}
