package sys_error

import (
	"fmt"
	"strconv"
)

const (
	SYS_ERR_PREFIX = 1000 //错误码为1000开头的为系统定义错误，其他微服务可以用其他错误码自定义错误

	RC_FAILURE                  = -1 //错误默认返回码
	RC_RENEGOTIATE_SESSION_KEY  = -2 //返回错误码:-2(客户端会触发重新协商会话密钥的动作)
	RC_RENEGOTIATE_REGISTER_KEY = -3 //返回错误码:-3(客户端会触发重新协商注册密钥的动作)
)

// 错误码格式为：错误码前缀(用于区分不同微服务)+内部错误码(系统内各个模块定义的错误码)
func errorCode(code int64) int64 {
	res, _ := strconv.ParseInt(fmt.Sprintf("%d%d", SYS_ERR_PREFIX, code), 10, 64)
	return res
}

var (
	//Authenticate
	AuthorizeFailure  = SysError{errorCode(100001), "Authorize failure", RC_FAILURE}
	TokenNotFound     = SysError{errorCode(100002), "Token not found", RC_FAILURE}
	ChannelIdNotFound = SysError{errorCode(100003), "Channel id is not found", RC_FAILURE}
	AuthenticateError = SysError{errorCode(100004), "The request user info is incorrect", RC_FAILURE}

	//授权错误
	SessionNotFoundError = SysError{errorCode(101001), "This equipment session was not found", RC_FAILURE}

	//数据验证
	RegisterKeyError           = SysError{errorCode(102001), "Register key error", RC_RENEGOTIATE_REGISTER_KEY}
	SessionKeyError            = SysError{errorCode(102002), "Session key error", RC_RENEGOTIATE_SESSION_KEY}
	HMACNotFound               = SysError{errorCode(102003), "HMAC not found", RC_FAILURE}
	UnknownCryptType           = SysError{errorCode(102004), "Unknown crypt type", RC_FAILURE}
	RegisterPostIntegrityError = SysError{errorCode(102005), "Register Post content is modified", RC_RENEGOTIATE_REGISTER_KEY}
	SessionPostIntegrityError  = SysError{errorCode(102006), "Session Post content is modified", RC_RENEGOTIATE_SESSION_KEY}
)

type SysError struct {
	Code    int64
	EMsg    string
	RetCode int32 //对外部返回的错误码（可选字段，其他服务可以不用添加）
}

const ERROR_SEPERATOR = "|||"

func (this SysError) Error() string {
	return fmt.Sprintf("%sret_code=%v%serr_code=%v%serr_msg=%s%s", ERROR_SEPERATOR, this.RetCode, ERROR_SEPERATOR, this.Code, ERROR_SEPERATOR, this.EMsg, ERROR_SEPERATOR)
}

// 可选方法，其他服务可以不用添加
func (this SysError) ReturnCode() int32 {
	return this.RetCode
}
