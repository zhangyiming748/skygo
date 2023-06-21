package response

import (
	"encoding/base64"
	"reflect"
	"strconv"
	// "go/src/cmd/vet/testdata/src/method"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"skygo_detection/guardian/app/http/request"
	pb "skygo_detection/guardian/app/http/response/proto"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/util"

	"skygo_detection/guardian/app/sys_error"
)

const (
	HTTP_RESPONSE_BODY = "Http-Response-Body"

	RC_SUCCESS = 0  // 成功默认返回码
	RC_FAILURE = -1 // 失败默认返回码
)

func GetResponseBody(ctx *gin.Context) string {
	if body, exist := ctx.Get(HTTP_RESPONSE_BODY); exist {
		return body.(string)
	} else {
		return ""
	}
}

func RenderSuccess(ctx *gin.Context, content *qmap.QM) {
	Render(ctx, RC_SUCCESS, "", content)
}

type ReturnErrorCode interface {
	ReturnCode() int32
}

func RenderFailure(ctx *gin.Context, err error) {
	errCode, errMsg := getRetErrCode(err)
	Render(ctx, errCode, errMsg, new(qmap.QM))
}

func getRetErrCode(err error) (int32, string) {
	var errCode int32 = RC_FAILURE
	var errMsg = err.Error()
	errInterface := reflect.TypeOf(new(ReturnErrorCode)).Elem()
	if reflect.TypeOf(err).Implements(errInterface) {
		// 如果自定义了错误返回码，则使用该错误返回码
		errCode = err.(ReturnErrorCode).ReturnCode()
	} else if idx := strings.Index(errMsg, "ret_code="); idx > -1 {
		// 否则尝试从错误信息里面取错误返回码
		// 如果错误格式是形如 rpc error: code = Unknown desc = ||| ret_code=-2 ||| err_code=1001101001 ||| err_msg=The auth token is expired |||
		// 则取其中的 ret_code 作为错误返回码
		if s := strings.Split(errMsg[idx:], sys_error.ERROR_SEPERATOR); len(s) > 0 {
			if codeArr := strings.Split(s[0], "="); len(codeArr) == 2 {
				if i, err := strconv.Atoi(codeArr[1]); err == nil {
					errCode = int32(i)
				}
			}
		}
	}

	// 从错误信息里面提取err_msg
	if idx := strings.Index(errMsg, "err_msg="); idx > -1 {
		if s := strings.Split(errMsg[idx:], sys_error.ERROR_SEPERATOR); len(s) > 0 {
			if msgArr := strings.Split(s[0], "="); len(msgArr) == 2 {
				errMsg = msgArr[1]
			}
		}
	} else {
		// 只提取错误信息
		errMsg = strings.Replace(errMsg, "rpc error: code = Unknown desc = ", "", 1)
	}
	return errCode, errMsg
}

func Render(ctx *gin.Context, code int32, msg string, content *qmap.QM) {
	if request.IsUseProtoBuffer(ctx) {
		renderProtoBuf(ctx, code, msg, content)
	} else {
		renderJson(ctx, code, msg, content)
	}
}

func getProtoBufBody(ctx *gin.Context, code int32, msg string, content *qmap.QM) []byte {
	dataStr := getResponseDataStr(content)
	var dataContent string
	if rsaKey := request.GetRsaMasterKey(ctx); rsaKey != "" {
		if dataStr == "" && sys_service.ENV != "online" {
			dataStr = request.GetRequestBody(ctx).ToString()
		}
		if dataStr != "" {
			encodedData, err := new(sys_service.CryptService).RsaEncrypt([]byte(dataStr), rsaKey)
			if err != nil {
				panic(err)
			}
			dataContent = base64.StdEncoding.EncodeToString(encodedData)
		}
	} else {
		if dataStr == "" && sys_service.ENV != "online" {
			dataContent = request.GetRequestBody(ctx).ToString()
		} else {
			dataContent = dataStr
		}
	}
	resp := &pb.Response{
		Code: code,
		Msg:  msg,
		Data: dataContent,
	}
	ctx.Set(HTTP_RESPONSE_BODY, resp.String())
	if binaryBody, err := proto.Marshal(resp); err == nil {
		return binaryBody
	} else {
		panic(err)
	}
}

func renderProtoBuf(ctx *gin.Context, code int32, msg string, data *qmap.QM) {
	finalBody := getProtoBufBody(ctx, code, msg, data)
	if aesKey := request.GetAesKey(ctx); aesKey != "" {
		keyByte := []byte(aesKey)
		if encodedResponse, err := new(sys_service.CryptService).AesEncrypt(finalBody, keyByte[0:16], keyByte[16:32]); err == nil {
			finalBody = encodedResponse
			mac := util.CalcHMACMd5(keyByte[0:16], finalBody)
			ctx.Header("X-HMAC", mac)
		} else {
			panic(err)
		}
	}

	ctx.Abort()
	ctx.Data(200, ".*", finalBody)
}

func renderJson(ctx *gin.Context, code int32, msg string, content *qmap.QM) {
	finalBody := qmap.QM{"code": code}
	if rsaKey := request.GetRsaMasterKey(ctx); rsaKey != "" {
		if dataStr := getResponseDataStr(content); dataStr != "" {
			encodedData, err := new(sys_service.CryptService).RsaEncrypt([]byte(dataStr), rsaKey)
			if err != nil {
				panic(err)
			}
			finalBody["data"] = base64.StdEncoding.EncodeToString(encodedData)
		}
	} else {
		if content != nil && *content != nil {
			data := getResponseData(content)
			switch data.(type) {
			case string:
				if data.(string) != "" {
					finalBody["data"] = data
				}
			case qmap.QM:
				if len(data.(qmap.QM)) > 0 {
					finalBody["data"] = data
				}
			default:
				finalBody["data"] = data
			}
		}
	}
	if msg != "" {
		finalBody["msg"] = msg
	}
	ctx.Set(HTTP_RESPONSE_BODY, finalBody.ToString())
	ctx.AbortWithStatusJSON(200, finalBody)
}

func getResponseData(resp *qmap.QM) interface{} {
	if val, exist := (*resp)["data"]; exist {
		return val
	}
	return resp
}

func getResponseDataStr(resp *qmap.QM) (dataStr string) {
	data := util.ResolvePointValue(getResponseData(resp))
	switch reflect.TypeOf(data).Kind() {
	case reflect.String:
		dataStr = data.(string)
	case reflect.Slice:
		dataStr = util.SliceToString(data.([]interface{}))
	default:
		switch data.(type) {
		case qmap.QM:
			dataStr = data.(qmap.QM).ToString()
		default:
			dataStr = util.MapToString(data.(map[string]interface{}))
		}
	}
	return
}
