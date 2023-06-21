package middleware

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/request"
)

// 对Body进行JSON解析
func JsonDecode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//如果body还没有被解码， 则尝试使用JSON进行解析
		if !request.IsSetBody(ctx) {
			decodeJsonBody(ctx)
		}
		ctx.Next()
	}
}

// 解码JSON传输的body
func decodeJsonBody(ctx *gin.Context) {
	//如果请求内容类型不是multipart/form-data，则默认使用json进行解析
	//注：此处之所以对multipart/form-data进行特别处理，是因为ioutil.ReadAll(ctx.Request.Body)函数读完之后，ctx.Request.Body内的数据被清空
	//会导致无法取出表单里面的数据，如果能解决此问题，则不用对multipart/form-data进行特别处理
	//该问题后续研究一下
	contentType := strings.ToLower(ctx.ContentType())
	if contentType != "multipart/form-data" {
		bodyBinary, _ := ioutil.ReadAll(ctx.Request.Body)
		body := getDecodeBody(ctx, bodyBinary)
		request.SetBody(ctx, body)
	}
}

func getDecodeBody(ctx *gin.Context, body []byte) *qmap.QM {
	if len(body) > 0 {
		requestBody := qmap.QM{}
		json.Unmarshal(body, &requestBody)
		return &requestBody
	} else {
		return new(qmap.QM)
	}
}
