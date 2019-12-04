package medusa

import (
	"testing"
)

func TestString(t *testing.T) {
	s := "hellow world!"
	sd := encode([]byte(s))
	s2 := decode(sd)
	if s != string(s2) {
		t.Fatal(s, "!=", string(s2), "fail to encode and decode")
	}
}
func TestBytes(t *testing.T) {
	s := []byte{5, 1, 0}
	sd := encode([]byte(s))
	s2 := decode(sd)
	for i := range s {
		if s[i] != s2[i] {
			t.Fatal(s, "!=", string(s2), "fail to encode and decode")
			return
		}
	}
}
