package controller

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	mycaptcha "skygo_detection/lib/common_lib/captcha"
)

type CollectionController struct{}

/**
 * apiType http
 * @api {get} /api/v1/captcha 获取验证码图片
 * @apiVersion 1.0.0
 * @apiName GetCaptcha
 * @apiGroup Route
 *
 * @apiDescription 获取验证码图片
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/captcha
 */
func (this CollectionController) GetCaption(ctx *gin.Context) {
	//生成验证码的id，并把id做为redis的key，具体内容作为redis的value
	captchaId := mycaptcha.NewLen(ctx, 4)

	//渲染出图片
	if Serve(ctx.Writer, ctx.Request, captchaId, ".png", "en", false, captcha.StdWidth, captcha.StdHeight) == captcha.ErrNotFound {
		http.NotFound(ctx.Writer, ctx.Request)
	}
}

//验证
//func (this CollectionController) Verify(c *gin.Context) {
//	captchaId, err := c.Request.Cookie("captcha_")
//	if err != nil {
//		panic(err)
//	}
//	value := c.Param("value")
//	if captchaId.Value == "" || value == "" {
//		c.String(http.StatusBadRequest, "参数错误")
//	}
//	fmt.Println(captchaId, "  ", value)
//	if captcha.VerifyString(captchaId.Value, value) {
//		c.JSON(http.StatusOK, "验证成功")
//	} else {
//		c.JSON(http.StatusOK, "验证失败")
//	}
//}

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
