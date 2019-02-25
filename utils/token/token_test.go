package token

import (
	"crypto/md5"
	"testing"
	"time"

	"crypto/hmac"
	"go_base/utils/encoding"
	"go_base/utils/testing2"
)

func TestCipher(t *testing.T) {
	tt := testing2.Wrap(t)
	c := NewCipher([]byte("12345"), time.Second*100, md5.New, encoding.Base64URL)

	for _, s := range []string{"a", "b", "c", "d", "e", "fdddddddddddddddddddddd"} {
		tok := c.Encode([]byte(s))
		tt.Log(string(tok), len(tok))

		ds, err := c.Decode(tok)
		tt.DeepEq([]byte(s), ds).Nil(err)
	}

	for _, s := range []string{"a", "b", "c", "d", ""} {
		tok := c.Encode([]byte(s))
		ds, err := c.Decode(tok)
		tt.DeepEq([]byte(s), ds).Nil(err)
	}
}

var cipher = NewCipher([]byte("12345"), time.Second*100, md5.New)
var data = []byte("abcdefghijklmn")
var encData = cipher.Encode(data)

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cipher.Encode(data)
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cipher.Decode(encData)
	}
}

func BenchmarkHMac(b *testing.B) {
	key := []byte("12345")
	for i := 0; i < b.N; i++ {
		hm := hmac.New(md5.New, key)
		hm.Write(data)
		_ = hm.Sum(nil)
	}
}
