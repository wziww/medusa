package padding

import "testing"

func TestStringZero(t *testing.T) {
	so := "hello world"
	sp := ZeroPadding([]byte(so), 16)
	sup, _ := ZeroUnPadding(sp, 16)
	if so != string(sup) {
		t.Fatal(so, "!=", string(sup), "fail to ZERO pad and unpad")
	}
}

func TestZero(t *testing.T) {
	bo := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	bp := ZeroPadding([]byte(bo), 16)
	bup, _ := ZeroUnPadding(bp, 16)
	for i := range bo {
		if bo[i] != bup[i] {
			t.Fatal(bo, "!=", bup, "fail to ZERO pad and unpad")
			return
		}
	}
}
