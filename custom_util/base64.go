package custom_util

import "encoding/base64"

// base64解密
func Base64Decode(s string) (string, error) {
	sBytes, err := base64.StdEncoding.DecodeString(s)
	return string(sBytes), err
}

// base64加密
func Base64Encode(s string) string {
	sBytes := []byte(s)
	return base64.StdEncoding.EncodeToString(sBytes)
}
