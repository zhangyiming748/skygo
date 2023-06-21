package custom_util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 算法 Aes
// 模式 CBC
// 补码 PKCS5Padding
// iv 初始向量

// Aes ...
type Aes struct {
	Key []byte
	Iv  []byte
}

// Encrypt 加密
func (t *Aes) Encrypt(b []byte) ([]byte, error) {
	block, err := aes.NewCipher(t.Key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	b = t.pkcs5Padding(b, blockSize)

	iv := make([]byte, blockSize)

	if len(t.Iv) != 0 && len(t.Iv) != blockSize {
		return nil, errors.New("Aes iv length err")
	}

	if len(t.Iv) == 0 {
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
	} else {
		iv = t.Iv
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)

	crypted := make([]byte, len(b))
	blockMode.CryptBlocks(crypted, b)

	return crypted, nil
}

// Decrypt 解密
func (t *Aes) Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(t.Key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, t.Key[:blockSize])

	b := make([]byte, len(crypted))
	blockMode.CryptBlocks(b, crypted)

	return t.pkcs5UnPadding(b), nil
}

// Base64Encode ...
func (t *Aes) Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base64Decode ...
func (t *Aes) Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func (t *Aes) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (t *Aes) pkcs5UnPadding(b []byte) []byte {
	length := len(b)
	unpadding := int(b[length-1])
	return b[:(length - unpadding)]
}
