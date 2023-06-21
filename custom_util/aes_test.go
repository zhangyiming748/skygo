package custom_util

import (
	"encoding/base64"
	"testing"
)

func TestEncrypt(t *testing.T) {
	aes := Aes{
		Key: []byte("1234567890abcDEF"),
		Iv:  []byte("1234567890abcDEF"),
	}

	rst, err := aes.Encrypt([]byte("1"))

	t.Log(base64.StdEncoding.EncodeToString(rst))
	t.Log(err)

	t.Fatal("end ...")
}

func TestDecrypt(t *testing.T) {
	aes := Aes{
		Key: []byte("1234567890abcDEF"),
	}

	dec, _ := aes.Base64Decode("Jy7FHsGpoYBC/uawk0jX1A==")
	b, err := aes.Decrypt(dec)

	t.Log(string(b))
	t.Log(err)

	t.Fatal("end ...")
}
