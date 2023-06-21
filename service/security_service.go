package service

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var defaultHashKey = "$1sQLU&*A=BxuGF%"
var defaultCost = 10
var extraHMD5String = "y8&NS1)^1`JL/"

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

func CheckPassword(hashedPassword, password string) error {
	fmt.Println(hashedPassword, password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// AccountAuthenticateError
		return errors.New("账户密码错误")
	}
	return nil
}

func HMD5(content, key string) string {
	hashKey := defaultHashKey
	if len(key) > 0 {
		hashKey = key
	}
	hMacs := hmac.New(md5.New, []byte(hashKey))
	hMacs.Write([]byte(content + extraHMD5String))
	return base64.StdEncoding.EncodeToString(hMacs.Sum([]byte(nil)))
}

func CheckHMD5(content, key, hashCode string) bool {
	if HMD5(content, key) == hashCode {
		return true
	}
	return false
}
