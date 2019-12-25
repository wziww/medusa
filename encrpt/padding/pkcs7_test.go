package padding

import (
	"testing"
)

func TestStringPKCS7(t *testing.T) {
	so := "hello world"
	sp := PKCS7Padding([]byte(so), 16)
	sup, _ := PKCS7UnPadding(sp, 16)
	if so != string(sup) {
		t.Fatal(so, "!=", string(sup), "fail to PKCS7 pad and unpad")
	}
}

func TestBytesPKCS7(t *testing.T) {
	bo := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	bp := PKCS7Padding([]byte(bo), 16)
	bup, _ := PKCS7UnPadding(bp, 16)
	for i := range bo {
		if bo[i] != bup[i] {
			t.Fatal(bo, "!=", bup, "fail to PKCS7 pad and unpad")
			return
		}
	}
}
