package padding

import (
	"testing"
)


func TestStringISO10126(t *testing.T) {
	so := "hello world"
	sp := ISO10126Padding([]byte(so), 16)
	sup, _ := ISO10126UnPadding(sp, 16)
	if so != string(sup) {
		t.Fatal(so, "!=", string(sup), "fail to ISO10126 pad and unpad")
	}
}


func TestBytesISO10126(t *testing.T) {
	bo := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	bp := ISO10126Padding([]byte(bo), 16)
	bup, _ := ISO10126UnPadding(bp, 16)
	for i := range bo {
		if bo[i] != bup[i] {
			t.Fatal(bo, "!=", bup, "fail to ISO10126 pad and unpad")
			return
		}
	}
}
