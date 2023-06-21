package session

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"skygo_detection/guardian/app/sys_error"
	"skygo_detection/guardian/src/net/metadata"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/util"
	"strconv"
	"time"
)

const (
	// Request
	SOURCE_SERVICE      = "source_service"      //请求来源服务名称
	SESSION_ID          = "session_id"          //session id
	REQUEST_ID          = "request_id"          //请求全局id（通过全局id可以查询到参与某次响应服务的所有节点的服务日志）
	REQUEST_INNER_ID    = "request_inner_id"    //请求内部id（通过内部id可以查询到单次响应服务的所有节点的服务拓扑顺序）
	REQUEST_INNER_COUNT = "request_inner_count" //请求内部服务调用计数（单次服务，单个节点的服务调用计数，用于生成具有拓扑顺序关系的请求内部id）
	LOG_OUTPUT_FLAG     = "log_output_flag"     //日志输出标记

	// User Info
	USER_INFO_NAME       = "username"
	USER_INFO_ID         = "user_id"
	USER_ROLE_ID         = "role_id"
	USER_GLOBAL_ROLE_ID  = "global_role_id"
	USER_INFO_CHANNEL_ID = "channel_id"
	USER_INFO_HMAC       = "hmd5"
	USER_ACCOUNT_TYPE    = "accout_type" //用户账户类型(user、vehicle、tbox)

	ADMIN_CHANNEL_ID = "Q00001" //奇虎管理员渠道号
)

// 需要传输到其他服务的上下文关键字
var outputCtxKey = map[string]interface{}{
	SESSION_ID:           struct{}{},
	REQUEST_ID:           struct{}{},
	USER_INFO_NAME:       struct{}{},
	USER_INFO_ID:         struct{}{},
	USER_INFO_CHANNEL_ID: struct{}{},
	USER_INFO_HMAC:       struct{}{},
	USER_GLOBAL_ROLE_ID:  struct{}{},
	LOG_OUTPUT_FLAG:      struct{}{},
}

// 获取输出到其他服务的上下文
func NewOutputContext(ctx context.Context) context.Context {
	newMD := qmap.QM{}
	if rawMD, ok := metadata.FromContext(ctx); ok {
		for key, val := range rawMD {
			if _, ok := outputCtxKey[key]; ok {
				newMD[key] = val
			}
		}
	}
	newMD[REQUEST_INNER_ID] = NewServiceInvokeId(ctx)

	return metadata.NewContext(context.Background(), newMD)
}

func Get(ctx context.Context, key string) interface{} {
	return metadata.Value(ctx, key)
}

func GetString(ctx context.Context, key string) string {
	return metadata.String(ctx, key)
}

func Set(ctx context.Context, key string, value interface{}) {
	if md, ok := metadata.FromContext(ctx); ok {
		md[key] = value
	}
}

func MustGet(ctx context.Context, key string) string {
	if md, ok := metadata.FromContext(ctx); ok {
		if val, exist := md[key]; exist {
			return val.(string)
		}
	}
	panic(sys_error.SessionNotFoundError)
}

func GetSessionId(ctx context.Context) string {
	return metadata.String(ctx, SESSION_ID)
}

func GetChannelId(ctx context.Context) string {
	return metadata.String(ctx, USER_INFO_CHANNEL_ID)
}

func GetAccountType(ctx context.Context) string {
	return metadata.String(ctx, USER_ACCOUNT_TYPE)
}

// 获取查询渠道号
// 如果渠道号是奇虎管理员渠道号，则直接返回空字符串
func GetQueryChannelId(ctx context.Context) string {
	channelId := metadata.String(ctx, USER_INFO_CHANNEL_ID)
	if channelId != ADMIN_CHANNEL_ID {
		return channelId
	}
	return ""
}

func GetUserId(ctx context.Context) int64 {
	return metadata.Int64(ctx, USER_INFO_ID)
}

func GetGlobalRoleId(ctx context.Context) int64 {
	return metadata.Int64(ctx, USER_GLOBAL_ROLE_ID)
}

func GetRoleId(ctx context.Context) int64 {
	return metadata.Int64(ctx, USER_ROLE_ID)
}

func GetHMAC(ctx context.Context) string {
	return metadata.String(ctx, USER_INFO_HMAC)
}

func GetUserName(ctx context.Context) string {
	return metadata.String(ctx, USER_INFO_NAME)
}

func GetSourceService(ctx context.Context) string {
	return metadata.String(ctx, SOURCE_SERVICE)
}

// 获取请求全局唯一id
func GetRequestId(ctx context.Context) string {
	requestId := metadata.String(ctx, REQUEST_ID)
	if requestId == "" {
		requestId = GenerateRequestId()
		Set(ctx, REQUEST_ID, requestId)
	}
	return requestId
}

// 获取当前请求内部id
func GetRequestInnerId(ctx context.Context) string {
	requestInnerId := metadata.String(ctx, REQUEST_INNER_ID)
	if requestInnerId == "" {
		requestInnerId = "1"
		Set(ctx, requestInnerId, REQUEST_INNER_ID)
	}
	return requestInnerId
}

// 生成新的服务调用id,该id将作为被调用服务的内部id
func NewServiceInvokeId(ctx context.Context) string {
	innerId := GetRequestInnerId(ctx)
	addedInnerCount := addRequestInnerCount(ctx)
	return fmt.Sprintf("%s.%d", innerId, addedInnerCount)
}

// 增加服务内部调用次数,并返回增加后的数值
func addRequestInnerCount(ctx context.Context) int64 {
	innerCount := metadata.Int64(ctx, REQUEST_INNER_COUNT)
	innerCount++
	Set(ctx, REQUEST_INNER_COUNT, innerCount)
	return innerCount
}

// 生成请求全局唯一id
func GenerateRequestId() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%v-%v-%v", hostname, util.GetCurrentMilliSecond(), rand.Intn(1000000))
}

func GenerateSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// 获取日志flag
func GetLogOutputFlag(ctx context.Context) bool {
	return metadata.Bool(ctx, LOG_OUTPUT_FLAG)
}

// 设置日志flag
func SetLogOutputFlag(ctx context.Context, probability float32) bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(10000) <= int(probability*10000) {
		Set(ctx, LOG_OUTPUT_FLAG, true)
	} else {
		Set(ctx, LOG_OUTPUT_FLAG, false)
	}
	return true
}
