package captcha

import (
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"skygo_detection/lib/common_lib/redis"
)

var prefix = "skygo_detection:captcha:"

type RedisStorage struct {
	expiration time.Duration
}

// Set sets the digits for the captcha id.
func (t *RedisStorage) Set(id string, digits []byte) {
	redis.NewRedis(0).Set(prefix+id, digits, t.expiration)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (t *RedisStorage) Get(id string, clear bool) (digits []byte) {
	// 不处理clear情况

	strC := redis.NewRedis(0).Get(prefix + id)
	return []byte(strC.Val())
}

var redisStorage = &RedisStorage{time.Second * 60}

// --------------------------------------------
func init() {
	captcha.SetCustomStore(redisStorage)
}

func NewLen(ctx *gin.Context, i int) string {
	captchaId := captcha.NewLen(i)
	// 写到cookie中
	uid_cookie := &http.Cookie{
		Name:     "captcha_",
		Value:    captchaId,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   120,
	}
	http.SetCookie(ctx.Writer, uid_cookie)
	return captchaId
}

func Verify(ctx *gin.Context, captchaText string) bool {
	captchaId, err := ctx.Request.Cookie("captcha_")
	if err != nil {
		panic(err)
	}
	if captchaId.Value == "" || captchaText == "" {
		return false
	}
	defer redis.NewRedis(0).Del(prefix + captchaId.Value)
	if captcha.VerifyString(captchaId.Value, captchaText) {
		return true
	} else {
		return false
	}
}
