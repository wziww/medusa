package encrypt

import (
	"testing"
)

var passwordCfb []byte = []byte("iT87ok9Uhn")

var obj = (&CipherS{
	Password:    &passwordCfb,
	Method:      "aes-256-cfb",
	PaddingMode: "",
}).Construct("aes-256-cfb", []byte("0000000000000000")).(*CipherS)
var errobj = (&CipherS{
	Password:    &passwordCfb,
	Method:      "asd2e",
	PaddingMode: "",
})

func TestConstruct(t *testing.T) {
	result := errobj.Construct("aes-256-cfb", []byte("0000000000000000"))
	errResult := errobj.Construct("xv34rt2345", []byte("0000000000000000"))
	if result == nil || errResult != nil {
		t.Fatal("construct error")
	}
}
func TestStringCfb(t *testing.T) {
	s := `abc`
	obj.encodeiv = []byte("0000000000000000")
	sd := obj.Encode([]byte(s))
	s2 := obj.Decode(sd)
	if s != string(s2) {
		t.Fatal(s, "!=", string(s2), "fail to encode and decode")
	}
}

func benchmarkEncode(b *testing.B, buf []byte) {
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj.Encode([]byte(buf))
	}
}

func benchmarkDecode(b *testing.B, buf []byte) {
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj.Decode(buf)
	}
}
