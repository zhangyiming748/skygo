package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	mrand "math/rand"
)

type Crypt_service struct{}

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8z1GclvVwxJ854TGCl2dsHazJgIdRxFADf5oKjnbFwF2ei/U40m0nvgObZLnvEfLaUfdgbczmcAIocOIVHzec1wMKbNNDjxv9LH9npyXrz83IsQSq81qJf7+q/5AlZ0TBg+bF3eKztQUi+Hpg9k7O5Lkv/uGAfRGGvDO35zBc6wIDAQAB
-----END PUBLIC KEY-----`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBALzPUZyW9XDEnznhMYKXZ2wdrMmAh1HEUAN/mgqOdsXAXZ6L9TjSbSe+A5tkue8R8tpR92BtzOZwAihw4hUfN5zXAwps00OPG/0sf2enJevPzcixBKrzWol/v6r/kCVnRMGD5sXd4rO1BSL4emD2Ts7kuS/+4YB9EYa8M7fnMFzrAgMBAAECgYBZEt2HqFgmWTxdC/ZVi6QJB37qmS49zwWIgPxlGozCAlyoXZLUucExTJ1bBAwL00Xk5WJ1JZfS5ui9t3ORT2bmUAakxpeU9UsrrC1B9Um60hgtVDWd1eXIYL5IlM+Aqvk3FHRDF2PrDVYCZLavKHNd0l/msrv5tZ7zFfsuGQUywQJBAPobP2Rm19XEBq868KFZj99uzxni2iEeDnpxYt+ijRKeW+aTEfvflfU6I/IUW2DgLrPEGgIrFJqU7XdP3UwPWtECQQDBQk2UyFUHlpWivCA6PwTq/1AhHkFEZ1Gd1XdPvGZZTeZ0HV+QA/1/IlohuuR80dnBdi2hJarZUthjSrk8bbL7AkEA5YGZc124U84VQDl61OUl1CeP3jZAakF1kcB4tbUpdVtiA70TtKjgp+6ZS6yIieZOlOGv6Ct2Nb/SBTmBXil88QJAVKcZWpmh/U/tvbnQGBNwsQsi607YYgEr1AokWA37exTPZH9VU70btiuy9WFrIm29h6ufcx4Px2AtntilaR3YLwJBANrHz0PutDyvYuHLuD961pRPS6KMv6ndaNFGIAmWYLKoZRHTKpH+uGGoBMYcAmv18isuOxlx3gmA34ov/WQbAVA=
-----END PRIVATE KEY-----`)
)

var (
	pubRSAKeyInstance *rsa.PublicKey
	priRSAKeyInstance *rsa.PrivateKey
)

func (c *Crypt_service) DecodeRequest(ctx *gin.Context) []byte {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if masterKey := c.GetRequstMaterKey(ctx); len(body) > 0 && masterKey != "" {
		key, iv := c.GetKeyIV(masterKey)
		if aesDecoded, err := c.AesDecrypt(body, key, iv); err == nil {
			return aesDecoded
		} else {
			panic(err)
		}
	} else {
		return body
	}
}

func (c *Crypt_service) EncodeResponse(originData []byte, ctx *gin.Context) []byte {
	if masterKey := c.GetRequstMaterKey(ctx); masterKey != "" {
		key, iv := c.GetKeyIV(masterKey)
		if aesEncoded, err := c.AesEncrypt(originData, key, iv); err == nil {
			return aesEncoded
		} else {
			panic(err)
		}
	} else {
		return originData
	}
}

func (c *Crypt_service) GetKeyIV(key string) ([]byte, []byte) {
	if binayKey, err := base64.StdEncoding.DecodeString(key); err == nil {
		masterKey := string(c.RsaDecrypt(binayKey))
		return []byte(masterKey[17:49]), []byte(masterKey[49:65])
	} else {
		panic("cannot base64 decode master key ")
	}
}

func (c *Crypt_service) GetRequstMaterKey(ctx *gin.Context) (mk string) {
	mk = ctx.GetHeader("X-SkyGo-MK")
	if mk == "" {
		mk = ctx.Query("X-SKYGO-MK")
	}
	return
}

// RSA加密
func (c *Crypt_service) RsaEncrypt(origData []byte) []byte {
	pubKeyInstance := c.getRSAPubKeyInstance()
	if encoded, err := rsa.EncryptPKCS1v15(rand.Reader, pubKeyInstance, origData); err == nil {
		return encoded
	} else {
		panic(err)
	}
}

// RSA解密
func (c *Crypt_service) RsaDecrypt(cipherText []byte) []byte {
	priKeyInstance := c.getRSAPriKeyInstance()
	if decoded, err := rsa.DecryptPKCS1v15(rand.Reader, priKeyInstance, cipherText); err == nil {
		return decoded
	} else {
		panic(err)
	}
}

func (c *Crypt_service) getRSAPubKeyInstance() *rsa.PublicKey {
	if pubRSAKeyInstance == nil {
		pubKey, _ := pem.Decode(RsaPublicKey) //将密钥解析成公钥实例
		if pubKey == nil {
			panic("RSA public key error")
		}
		if keyInstance, err := x509.ParsePKIXPublicKey(pubKey.Bytes); err == nil {
			pubRSAKeyInstance = keyInstance.(*rsa.PublicKey)
		} else {
			panic(err)
		}
	}
	return pubRSAKeyInstance
}

func (c *Crypt_service) getRSAPriKeyInstance() *rsa.PrivateKey {
	if priRSAKeyInstance == nil {
		priKey, _ := pem.Decode(RsaPrivateKey) //将密钥解析成私钥实例
		if priKey == nil {
			panic("RSA private key error")
		}
		if keyInstance, err := x509.ParsePKCS8PrivateKey(priKey.Bytes); err == nil {
			priRSAKeyInstance = keyInstance.(*rsa.PrivateKey)
		} else {
			panic(err)
		}
	}
	return priRSAKeyInstance
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (c *Crypt_service) AesEncrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (c *Crypt_service) AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func (c *Crypt_service) GenerateRegisterKey() string {
	bytesArr := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()+-")
	masterKey := []byte{}

	for i := 0; i < 32; i++ {
		masterKey = append(masterKey, bytesArr[mrand.Intn(len(bytesArr))])
	}
	return string(masterKey)
}

func (c *Crypt_service) GenerateSessionKey() string {
	bytesArr := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()+-")
	masterKey := []byte{}

	for i := 0; i < 24; i++ {
		masterKey = append(masterKey, bytesArr[mrand.Intn(len(bytesArr))])
	}
	return string(masterKey)
}
