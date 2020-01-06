package  padding

import (
	"testing"
)

func TestStringANSIX923(t *testing.T){
so := "hello world"
	sp := ANSIX923Padding([]byte(so), 16)
	sup, _ := ANSIX923UnPadding(sp, 16)
	if so != string(sup) {
		t.Fatal(so, "!=", string(sup), "fail to ANSIX923 pad and unpad")
	}
}

func TestBytesANSIX923(t *testing.T){
	bo := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	bp := ANSIX923Padding([]byte(bo), 16)
	bup, _ := ANSIX923UnPadding(bp, 16)
	for i := range bo {
		if bo[i] != bup[i] {
			t.Fatal(bo, "!=", bup, "fail to ANSIX923 pad and unpad")
			return
		}
	}
}