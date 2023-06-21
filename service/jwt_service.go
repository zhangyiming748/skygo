package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"

	"strings"
	"time"
)

/**
 * 生成JWT
 * @params groupId 	string 所属组id
 * @params id 		string 用户id
 * @params subject 	string 主题
 */
func GenerateJWT(id, subject string, authorizeTime int) (string, error) {
	jwtConfig := LoadJWTConfig()
	secretKey := []byte(jwtConfig.SecretKey)
	var claims *JWTClaims
	if authorizeTime > 0 {
		claims = &JWTClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + int64(authorizeTime),
				Id:        id,
				IssuedAt:  time.Now().Unix(),
				Subject:   subject,
			},
		}
	} else {
		claims = &JWTClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + int64(jwtConfig.ExpireTime),
				Id:        id,
				IssuedAt:  time.Now().Unix(),
				Subject:   subject,
			},
		}
	}

	token := jwt.NewWithClaims(getSigningMethod(jwtConfig.Algorithm), claims)
	if jwtToken, err := token.SignedString(secretKey); err == nil {
		return jwtToken, nil
	} else {
		return "", err
	}
}

func getSigningMethod(algorithm string) *jwt.SigningMethodHMAC {
	switch strings.ToUpper(algorithm) {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		panic("unknown jwt sign algorithm")
	}
}

func TokenValid(token string) (*JWTClaims, error) {
	jwtConfig := LoadJWTConfig()
	secretKey := []byte(jwtConfig.SecretKey)
	jwtClaim := new(JWTClaims)
	if _, err := jwt.ParseWithClaims(token, jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}); err != nil {
		// TODO 错误统一定义
		return nil, errors.New("TokenInvalidError")
	}
	return jwtClaim, nil
}

// 自定义jwt payload结构体
// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type JWTClaims struct {
	jwt.StandardClaims
}

// Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c JWTClaims) Valid() error {
	now := jwt.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		return errors.New("TokenExpiredError")
	}

	if c.VerifyIssuedAt(now, false) == false {
		return errors.New("TokenUsedTooEarlyError")
	}

	if c.VerifyNotBefore(now, false) == false {
		return errors.New("TokenInvalidError")
	}

	return nil
}

// No errors
func (c *JWTClaims) validErr(e *jwt.ValidationError) bool {
	return e.Errors == 0
}
