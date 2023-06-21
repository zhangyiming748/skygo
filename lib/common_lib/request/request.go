package request

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
	"strings"
)

const (
	//预定义Header常量
	H_RSA_KEY      = "X-SkyGo-MK"    //RSA加密key
	H_CRYPT_TYPE   = "X-Crypt"       //AES加密类型 "1"：注册密钥加密 "2"：会话密钥加密
	H_HMAC         = "X-HMAC"        //HMAC校验值
	H_AUTH_TOKEN   = "Authorization" //授权令牌
	H_PROTO_BUFFER = "X-Proto-Buf"   //是否使用proto buffer进行数据传输(1:是, 其他:否)
	H_CHANNEL_ID   = "X-Channel-Id"  //渠道号
	H_SN           = "X-Sn"          //设备序列号

	HTTP_REQUEST_BODY = "Http-Request-Body"
	PROTO_BODY        = "Proto-Body"
	AES_KEY           = "Aes-Key"
)

const (
	EQUIPMENT_VEHICLE = "vehicle" //设备类型:车机
	EQUIPMENT_TBOX    = "tbox"    //设备类型:tbox

	ENCRYPT_RSA = "rsa" //加密算法:RSA
	ENCRYPT_AES = "aes" //加密算法:AES

	CRYPT_REGISTER = "1" //注册密钥加密
	CRYPT_SESSION  = "2" //会话密钥加密
)

func GetBasicAuthAccount(ctx *gin.Context) (string, string) {
	if rawAuth := ctx.GetHeader("Authorization"); rawAuth != "" {
		base64Account := strings.Replace(rawAuth, "Basic ", "", 1)
		if basic, err := base64.URLEncoding.DecodeString(base64Account); err == nil {
			if accountInfo := strings.Split(string(basic), ":"); len(accountInfo) == 2 {
				return accountInfo[0], accountInfo[1]
			}
		}
	}
	panic(errors.New("Auth info not found"))
}

func IsSetBody(ctx *gin.Context) bool {
	_, exist := ctx.Get(HTTP_REQUEST_BODY)
	return exist
}

func SetBody(ctx *gin.Context, body *qmap.QM) {
	ctx.Set(HTTP_REQUEST_BODY, body)
}

func GetRequestBody(ctx *gin.Context) *qmap.QM {
	if body, exist := ctx.Get(HTTP_REQUEST_BODY); exist {
		return body.(*qmap.QM)
	} else {
		return new(qmap.QM)
	}
}

func GetRequestQueryParams(ctx *gin.Context) *qmap.QM {
	u := url.URL{RawQuery: ctx.Request.URL.RawQuery}
	queryParam := qmap.QM{}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		queryParam[k] = v[0]
	}
	return &queryParam
}

func SetProtoBody(ctx *gin.Context, body []byte) {
	ctx.Set(PROTO_BODY, body)
}

func GetProtoBody(ctx *gin.Context) []byte {
	if body, exist := ctx.Get(PROTO_BODY); exist {
		return body.([]byte)
	} else {
		return []byte{}
	}
}

// 获取数据加密算法(RSA/AES)
func GetEncryptAlgorithm(ctx *gin.Context) string {
	if ctx.GetHeader(H_CRYPT_TYPE) != "" {
		return ENCRYPT_AES
	} else if ctx.GetHeader(H_RSA_KEY) != "" {
		return ENCRYPT_RSA
	}
	return ""
}

func GetRsaMasterKey(ctx *gin.Context) string {
	mk := ctx.GetHeader(H_RSA_KEY)
	if mk == "" {
		mk = ctx.Query("X-SKYGO-MK")
	}
	return mk
}

func SetAesKey(ctx *gin.Context, key string) {
	ctx.Set(AES_KEY, key)
}

func GetAesKey(ctx *gin.Context) string {
	if key, exist := ctx.Get(AES_KEY); exist {
		return key.(string)
	} else {
		return ""
	}
}

func GetHeaderSn(ctx *gin.Context) string {
	return ctx.GetHeader(H_SN)
}

func GetHMAC(ctx *gin.Context) string {
	return ctx.GetHeader(H_HMAC)
}

// AES加密类型 "1"：注册密钥加密 "2"：会话密钥加密
func GetCryptType(ctx *gin.Context) string {
	return ctx.GetHeader(H_CRYPT_TYPE)
}

func GetHeaderChannelID(ctx *gin.Context) string {
	return ctx.GetHeader(H_CHANNEL_ID)
}

// 从ctx中获取token内容
func GetAuthToken(ctx *gin.Context) string {
	token := ctx.GetHeader(H_AUTH_TOKEN)
	return strings.Replace(token, "Bearer ", "", -1)
}

// 是否使用proto buffer进行数据传输
func IsUseProtoBuffer(ctx *gin.Context) bool {
	if isUsePB := ctx.GetHeader(H_PROTO_BUFFER); isUsePB == "1" {
		return true
	} else {
		return false
	}
}

// 根据请求url判断请求的设备类型
func GetEquipmentType(ctx *gin.Context) string {
	pathArr := strings.Split(ctx.Request.URL.Path, "/")
	requestType, _ := SnakeString(pathArr[1])
	switch requestType {
	case "api":
		return EQUIPMENT_VEHICLE
	case "tbox_api":
		return EQUIPMENT_TBOX
	default:
		return ""
	}
}

func String(ctx *gin.Context, key string) (s string) {
	return GetRequestBody(ctx).String(key)
}

func DefaultString(ctx *gin.Context, key string, defaultVal string) (s string) {
	return GetRequestBody(ctx).DefaultString(key, defaultVal)
}

// must fetch a string ,or will panic
func MustString(ctx *gin.Context, key string) (s string) {
	return GetRequestBody(ctx).MustString(key)
}

func MustInt(ctx *gin.Context, key string) (s int) {
	return GetRequestBody(ctx).MustInt(key)
}

func Int(ctx *gin.Context, key string) (s int) {
	return GetRequestBody(ctx).Int(key)
}

func DefaultInt(ctx *gin.Context, key string, defaultVal int) (s int) {
	return GetRequestBody(ctx).DefaultInt(key, defaultVal)
}

func Int32(ctx *gin.Context, key string) (s int32) {
	return GetRequestBody(ctx).Int32(key)
}

func Int64(ctx *gin.Context, key string) (s int64) {
	return GetRequestBody(ctx).Int64(key)
}

func Float64(ctx *gin.Context, key string) (s float64) {
	return GetRequestBody(ctx).Float64(key)
}

func Float32(ctx *gin.Context, key string) (s float32) {
	return GetRequestBody(ctx).Float32(key)
}

func Bool(ctx *gin.Context, key string) (b bool) {
	return GetRequestBody(ctx).Bool(key)
}

func Slice(ctx *gin.Context, key string) (s []interface{}) {
	return GetRequestBody(ctx).Slice(key)
}

func MustSlice(ctx *gin.Context, key string) (s []interface{}) {
	return GetRequestBody(ctx).MustSlice(key)
}

func Map(ctx *gin.Context, key string) *qmap.QM {
	t := GetRequestBody(ctx).Map(key)
	return &t
}

func IsExist(ctx *gin.Context, key string) (b bool) {
	_, exist := GetRequestBody(ctx).TryInterface(key)
	return exist
}

/************************************/
/************ URL Query ************/
/************************************/
func QueryInt(ctx *gin.Context, key string) int {
	val := ctx.Query(key)
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return 0
}

func DefaultQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	if val, exist := ctx.GetQuery(key); exist {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func TryQueryInt(ctx *gin.Context, key string) (int, bool) {
	if val, exist := ctx.GetQuery(key); exist {
		if i, err := strconv.Atoi(val); err == nil {
			return i, true
		}
	}
	return 0, false
}

func QueryString(ctx *gin.Context, key string) string {
	return ctx.Query(key)
}

func MustQueryString(ctx *gin.Context, key string) string {
	if val, exist := ctx.GetQuery(key); exist {
		return val
	}
	panic("url params: " + key + " does not provided")
}

func DefaultQueryString(ctx *gin.Context, key string, defaultVal string) string {
	if val, exist := ctx.GetQuery(key); exist {
		return val
	}
	return defaultVal
}

/************************************/
/************ URL Param ************/
/************************************/
func ParamInt(ctx *gin.Context, key string) int {
	val := ctx.Param(key)
	if i, err := strconv.Atoi(val); err == nil {
		return i
	} else {
		panic(err)
	}
}

func ParamString(ctx *gin.Context, key string) string {
	return ctx.Param(key)
}

func SnakeString(s string) (string, bool) {
	data := make([]byte, 0, len(s)*2)
	change := false
	j := false
	pre := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			if i > 0 && j && pre {
				change = true
				data = append(data, '_')
			}
		} else {
			pre = true
		}

		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:])), change
}
