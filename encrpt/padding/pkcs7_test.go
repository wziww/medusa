package encrpt

import (
	// "fmt"
	"testing"
)

func TestStringPKCS7(t *testing.T) {
	so := "hello world"
	sp := PKCS7Padding([]byte(so), 16)
	// fmt.Printf("sp----%x\n", sp)
	sup, _ := PKCS7UnPadding(sp, 16)
	// fmt.Printf("sup----%x\n", sup)
	if so != string(sup) {
		t.Fatal(so, "!=", string(sup), "fail to PKCS7 pad and unpad")
	}
}

func TestBytesPKCS7(t *testing.T) {
	bo := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	bp := PKCS7Padding([]byte(bo), 16)
	// fmt.Printf("bp----%x\n", bp)
	bup, _ := PKCS7UnPadding(bp, 16)
	// fmt.Printf("sup----%x\n", bup)
	for i := range bo {
		if bo[i] != bup[i] {
			t.Fatal(bo, "!=", bup, "fail to PKCS7 pad and unpad")
			return
		}
	}
}
